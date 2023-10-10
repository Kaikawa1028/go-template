package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	stdErrors "errors"
	"github.com/Kaikawa1028/go-template/app/errors/types"
	"net/url"
	"os"

	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"gorm.io/gorm"
)

type Authenticate struct {
	db             *gorm.DB
	user           repository.User
	roleOptionRepo repository.AdminUserRoleOption
}

func NewAuthenticate(
	db *gorm.DB,
	user repository.User,
	roleOptionRepo repository.AdminUserRoleOption,
) *Authenticate {
	return &Authenticate{
		db:             db,
		user:           user,
		roleOptionRepo: roleOptionRepo,
	}
}

func (u Authenticate) Authenticate(requestToken string) (*model.User, error) {
	encKey := os.Getenv("ENC_KEY")
	if encKey == "" {
		return nil, errors.Wrap(errors.New("環境変数ENC_KEYが指定されていません"))
	}

	encIv := os.Getenv("ENC_IV")
	if encIv == "" {
		return nil, errors.Wrap(errors.New("環境変数ENC_IVが指定されていません"))
	}

	token, err := url.QueryUnescape(requestToken)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	key := []byte(encKey)
	ciphertext, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		if _, ok := err.(base64.CorruptInputError); ok {
			return nil, errors.Wrap(types.NewCantBase64DecodeAuthTokenError())
		}
		return nil, errors.Wrap(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	iv, err := base64.StdEncoding.DecodeString(encIv)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padSize := int(ciphertext[len(ciphertext)-1])
	ciphertext = ciphertext[:len(ciphertext)-padSize]
	hasher := sha256.New()
	hasher.Write(ciphertext)
	apiToken := hex.EncodeToString(hasher.Sum(nil))

	user, err := u.user.GetByApiToken(u.db, apiToken)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = u.attachRoleOption(user)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return user, nil
}

func (u Authenticate) AuthenticateByUserId(userId int) (*model.User, error) {
	user, err := u.user.GetByUserId(u.db, userId)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = u.attachRoleOption(user)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return user, nil
}

// attachRoleOption userに紐づくロールオプションを取得し、userにセットします。
// ロールオプションが存在しない場合はnilをセットします。
func (u Authenticate) attachRoleOption(user *model.User) error {
	roleOption, err := u.roleOptionRepo.GetAdminUserRoleOption(u.db, user.ID)
	if err != nil && !stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err)
	}

	user.RoleOption = roleOption

	return nil
}

package app

const (
	// RoleAdministrator システム管理者
	RoleAdministrator = "administrator"

	// RoleAgency 代理店
	RoleAgency = "agency"

	// RoleManager 管理者
	RoleManager = "manager"

	// RoleManagerCanFormatSetting 管理者：システム管理設定可能
	RoleManagerCanFormatSetting = "manager:can_format_setting"

	// RoleManagerCanSeoSetting 管理者：SEO情報管理可能
	RoleManagerCanSeoSetting = "manager:can_seo_setting"

	// RoleManagerCanTopSetting 管理者：店舗検索TOP管理可能
	RoleManagerCanTopSetting = "manager:can_top_setting"

	// RoleEditor 編集者
	RoleEditor = "editor"

	// RoleWorker 作業者
	RoleWorker = "worker"

	// RoleWorkerCanEditStore 作業者：店舗編集可能
	RoleWorkerCanEditStore = "worker:can_edit_store"

	// RoleViewer 閲覧者
	RoleViewer = "viewer"
)

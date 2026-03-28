package router

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&HostingAccountRouter{},
		&DashboardRouter{},
		&HostRouter{},
		&ContainerRouter{},
		&LogRouter{},
		&ToolboxRouter{},
		&CronjobRouter{},
		&BackupRouter{},
		&SettingRouter{},
		&AppRouter{},
		&WebsiteRouter{},
		&WebsiteDnsAccountRouter{},
		&WebsiteAcmeAccountRouter{},
		&WebsiteSSLRouter{},
		&DatabaseRouter{},
		&NginxRouter{},
		&RuntimeRouter{},
		&ProcessRouter{},
		&WebsiteCARouter{},
		&AIToolsRouter{},
		&GroupRouter{},
		&AlertRouter{},
	}
}

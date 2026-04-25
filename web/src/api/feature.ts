// COMMON (WEB AND SERVER)

export enum Feature {
    // Cluster features
    ViewClusterList = "view.cluster.list",
    ViewClusterItem = "view.cluster.item",
    ViewClusterOverview = "view.cluster.overview",
    ManageClusterCreate = "manage.cluster.create",
    ManageClusterUpdate = "manage.cluster.update",
    ManageClusterDelete = "manage.cluster.delete",

    // Tags features
    ViewTagList = "view.tag.list",

    // Node features
    ViewNodeDbOverview = "view.node.db.overview",
    ViewNodeDbConfig = "view.node.db.config",
    ViewNodeSshMetrics = "view.node.ssh.metrics",
    ManageNodeDbConfigUpdate = "manage.node.db.config.update",
    ManageNodeDbSwitchover = "manage.node.db.switchover",
    ManageNodeDbReinitialize = "manage.node.db.reinitialize",
    ManageNodeDbRestart = "manage.node.db.restart",
    ManageNodeDbReload = "manage.node.db.reload",
    ManageNodeDbFailover = "manage.node.db.failover",
    ManageNodeDbActivation = "manage.node.db.activation",
    ManageNodeSshDocker = "manage.node.ssh.docker",
    ViewNodeSshDocker = "view.node.ssh.docker",

    // Query features
    ViewQueryCrudList = "view.query.crud.list",
    ManageQueryCrudCreate = "manage.query.crud.create",
    ManageQueryCrudUpdate = "manage.query.crud.update",
    ManageQueryCrudDelete = "manage.query.crud.delete",
    ViewQueryDbInfo = "view.query.db.info",
    ViewQueryDbChart = "view.query.db.chart",
    ManageQueryDbTemplate = "manage.query.db.template",
    ManageQueryDbConsole = "manage.query.db.console",
    ManageQueryDbCancel = "manage.query.db.cancel",
    ManageQueryDbTerminate = "manage.query.db.terminate",
    ViewQueryLogList = "view.query.log.list",
    ManageQueryLogDelete = "manage.query.log.delete",

    // Vault features
    ViewVaultList = "view.vault.list",
    ManageVaultCreate = "manage.vault.create",
    ManageVaultUpdate = "manage.vault.update",
    ManageVaultDelete = "manage.vault.delete",

    // Cert features
    ViewCertList = "view.cert.list",
    ManageCertCreate = "manage.cert.create",
    ManageCertDelete = "manage.cert.delete",

    // Feature management features
    ViewPermissionList = "view.permission.list",
    ManagePermissionUpdate = "manage.permission.update",
    ManagePermissionDelete = "manage.permission.delete",

    // Bloat features
    ViewToolBloatList = "view.tool.bloat.list",
    ViewToolBloatItem = "view.tool.bloat.item",
    ViewToolBloatLogs = "view.tool.bloat.logs",
    ManageToolBloatJob = "manage.tool.bloat.job",

    // Management features
    ViewManagementSecret = "view.management.secret",
    ManageManagementSecret = "manage.management.secret",
    ManageManagementErase = "manage.management.erase",
    ManageManagementFree = "manage.management.free",
    ManageManagementBackup = "manage.management.backup",
}

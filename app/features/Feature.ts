// COMMON (WEB AND SERVER)

export enum Feature {
    // Cluster features
    ViewClusterList = "view.cluster.list",
    ViewClusterItem = "view.cluster.item", // NOTE: we do not use this api method in UI
    ViewClusterOverview = "view.cluster.overview",
    ManageClusterCreate = "manage.cluster.create",
    ManageClusterUpdate = "manage.cluster.update",
    ManageClusterDelete = "manage.cluster.delete",

    // Tags features
    ViewTagList = "view.tag.list",

    // Node features
    ViewNodeKeeperOverview = "view.node.keeper.overview", // NOTE: we do not use this api method in UI
    ViewNodeKeeperConfig = "view.node.keeper.config",
    ManageNodeKeeperConfigUpdate = "manage.node.keeper.config.update",
    ManageNodeKeeperSwitchover = "manage.node.keeper.switchover",
    ManageNodeKeeperReinitialize = "manage.node.keeper.reinitialize",
    ManageNodeKeeperRestart = "manage.node.keeper.restart",
    ManageNodeKeeperReload = "manage.node.keeper.reload",
    ManageNodeKeeperFailover = "manage.node.keeper.failover",
    ManageNodeKeeperActivation = "manage.node.keeper.activation",

    ViewNodePlatform = "view.node.platform",
    ManageNodePlatform = "manage.node.platform", // NOTE: we do not use this api method in UI
    ManageNodePlatformContainer = "manage.node.platform.container",
    ViewNodePlatformContainer = "view.node.platform.container",

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

    // Permission management features
    ViewPermissionList = "view.permission.list",
    ManagePermissionUpdate = "manage.permission.update",
    ManagePermissionDelete = "manage.permission.delete", // NOTE: we do not use this api method in UI

    // Tool features
    ViewToolPgCompactTableList = "view.tool.pg_compacttable.list",
    ViewToolPgCompactTableItem = "view.tool.pg_compacttable.item", // NOTE: we do not use this api method in UI
    ViewToolPgCompactTableLogs = "view.tool.pg_compacttable.logs",
    ManageToolPgCompactTableJob = "manage.tool.pg_compacttable.job",

    // Management features
    ViewManagementSecret = "view.management.secret", // NOTE: we do not use this api method in UI
    ManageManagementSecret = "manage.management.secret",
    ManageManagementErase = "manage.management.erase",
    ManageManagementFree = "manage.management.free",
    ManageManagementBackup = "manage.management.backup",
}

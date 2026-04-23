// COMMON (WEB AND SERVER)

export enum PermissionStatus {
    NOT_PERMITTED = 0,
    PENDING = 1,
    GRANTED = 2,
}

export enum Permission {
    // Cluster permissions
    ViewClusterList = "view.cluster.list",
    ViewClusterItem = "view.cluster.item",
    ViewClusterOverview = "view.cluster.overview",
    ManageClusterCreate = "manage.cluster.create",
    ManageClusterUpdate = "manage.cluster.update",
    ManageClusterDelete = "manage.cluster.delete",

    // Tags permissions
    ViewTagList = "view.tag.list",

    // Node permissions
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

    // Query permissions
    ViewQueryList = "view.query.list",
    ManageQueryCreate = "manage.query.create",
    ManageQueryUpdate = "manage.query.update",
    ManageQueryDelete = "manage.query.delete",
    ViewQueryExecuteInfo = "view.query.execute.info",
    ViewQueryExecuteChart = "view.query.execute.chart",
    ManageQueryExecuteTemplate = "manage.query.execute.template",
    ManageQueryExecuteConsole = "manage.query.execute.console",
    ManageQueryExecuteCancel = "manage.query.execute.cancel",
    ManageQueryExecuteTerminate = "manage.query.execute.terminate",
    ViewQueryLogList = "view.query.log.list",
    ManageQueryLogDelete = "manage.query.log.delete",

    // Vault permissions
    ViewVaultList = "view.vault.list",
    ManageVaultCreate = "manage.vault.create",
    ManageVaultUpdate = "manage.vault.update",
    ManageVaultDelete = "manage.vault.delete",

    // Cert permissions
    ViewCertList = "view.cert.list",
    ManageCertCreate = "manage.cert.create",
    ManageCertDelete = "manage.cert.delete",

    // Permission management
    ViewPermissionList = "view.permission.list",
    ManagePermissionUpdate = "manage.permission.update",
    ManagePermissionDelete = "manage.permission.delete",

    // Bloat permissions
    ViewBloatList = "view.bloat.list",
    ViewBloatItem = "view.bloat.item",
    ViewBloatLogs = "view.bloat.logs",
    ManageBloatJob = "manage.bloat.job",

    // Management permissions
    ViewManagementSecret = "view.management.secret",
    ManageManagementSecret = "manage.management.secret",
    ManageManagementErase = "manage.management.erase",
    ManageManagementFree = "manage.management.free",
    ManageManagementBackup = "manage.management.backup",
}

export type PermissionMap = {
    [key in Permission]: PermissionStatus;
};

export interface UserPermissions {
    username: string,
    permissions: PermissionMap,
}

// SPECIFIC (WEB)

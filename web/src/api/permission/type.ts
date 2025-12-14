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
    ManageClusterCreate = "manage.cluster.create",
    ManageClusterUpdate = "manage.cluster.update",
    ManageClusterDelete = "manage.cluster.delete",

    // Tags permissions
    ViewTagList = "view.tag.list",

    // Instance permissions
    ViewInstanceOverview = "view.instance.overview",
    ViewInstanceConfig = "view.instance.config",
    ManageInstanceConfigUpdate = "manage.instance.config.update",
    ManageInstanceSwitchover = "manage.instance.switchover",
    ManageInstanceReinitialize = "manage.instance.reinitialize",
    ManageInstanceRestart = "manage.instance.restart",
    ManageInstanceReload = "manage.instance.reload",
    ManageInstanceFailover = "manage.instance.failover",
    ManageInstanceActivation = "manage.instance.activation",

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

    // Password permissions
    ViewPasswordList = "view.password.list",
    ManagePasswordCreate = "manage.password.create",
    ManagePasswordUpdate = "manage.password.update",
    ManagePasswordDelete = "manage.password.delete",

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
}

export type PermissionMap = {
    [key in Permission]: PermissionStatus;
};

export interface UserPermissions {
    username: string,
    permissions: PermissionMap,
}

// SPECIFIC (WEB)

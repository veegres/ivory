import {ErrorSmart} from "../../../view/box/ErrorSmart";

export function ClusterNoInstanceError() {
    return <ErrorSmart error={"Main instance is not in the cluster, probably something has happened, you have some problems in your set up or you've change it"} />
}

export function ClusterNoLeaderError() {
    return <ErrorSmart error={"Main instance is not a leader, probably something has happened or you've change it"}/>
}

export function ClusterNoPostgresPassword() {
    return <ErrorSmart error={"You haven't set up postgres password for this cluster. Please, do it in the cluster settings bar"}/>
}

export function NoDatabaseError() {
    return <ErrorSmart error={"Database is not specified"}/>
}

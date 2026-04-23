import {ErrorSmart} from "../../../view/box/ErrorSmart"

export function ClusterNoNodeError() {
    return <ErrorSmart error={"Main node is not in the cluster, probably something has happened, you have some problems in your set up or you've change it"} />
}

export function ClusterNoLeaderError() {
    return <ErrorSmart error={"Main node is not a leader, probably something has happened or you've change it"}/>
}

export function ClusterNoPostgresVault() {
    return <ErrorSmart error={"You haven't set up postgres vault for this cluster. Please, do it in the cluster settings bar"}/>
}

export function NoDatabaseError() {
    return <ErrorSmart error={"Database is not specified"}/>
}

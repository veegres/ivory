import {ErrorAlert} from "../view/ErrorAlert";
import React from "react";

export function ClusterNoInstanceError() {
    return <ErrorAlert error={"Default instance is not in the cluster, probably something has happened, you have some problems in your set up or you've change it"} />
}

export function ClusterNoLeaderError() {
    return <ErrorAlert error={"Default instance is not a leader, probably something has happened or you've change it"}/>
}

export function ClusterNoPostgresPassword() {
    return <ErrorAlert error={"You haven't set up postgres password for this cluster. Please, do it in the cluster settings bar"}/>
}

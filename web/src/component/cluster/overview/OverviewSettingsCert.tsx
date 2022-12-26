import React, {useMemo} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {certApi, clusterApi} from "../../../app/api";
import {Autocomplete, Option} from "../../view/Autocomplete";
import {shortUuid} from "../../../app/utils";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {Cluster} from "../../../app/types";

type Props = {
    certId: string,
    cluster: Cluster,
}

export function OverviewSettingsCert(props: Props) {
    const { certId, cluster } = props

    const query = useQuery(["certs"], certApi.list)
    const options = useMemo(handleMemoOptions, [query.data])

    const updateMutationOptions = useMutationOptions(["cluster/list"])
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <Autocomplete
            label={"Patroni Cert"}
            selected={{key: certId, short: shortUuid(certId)}}
            options={options}
            loading={query.isLoading || updateCluster.isLoading}
            onUpdate={handleUpdate}
        />
    )

    function handleUpdate(option: Option | null) {
        updateCluster.mutate({...cluster, certId: option?.key})
    }

    function handleMemoOptions(): Option[] {
        return Object.entries(query.data ?? {}).map(([key, value]) => ({ key, short: shortUuid(key), name: value.fileName }))
    }
}

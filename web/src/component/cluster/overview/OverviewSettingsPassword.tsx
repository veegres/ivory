import {shortUuid} from "../../../app/utils";
import React, {useMemo} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {clusterApi, credentialApi} from "../../../app/api";
import {Cluster, CredentialType} from "../../../app/types";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {Autocomplete, Option} from "../../view/Autocomplete";

const keys = {
    [CredentialType.POSTGRES]: "postgresCredId",
    [CredentialType.PATRONI]: "patroniCredId",
}

type Props = {
    type: CredentialType
    cluster: Cluster
    credId: string
    label: string
}

export function OverviewSettingsPassword(props: Props) {
    const { type, label, cluster, credId } = props
    const passKey = keys[type]

    const query = useQuery(["credentials", type], () => credentialApi.list(type))
    const options = useMemo(handleMemoOptions, [query.data])

    const updateMutationOptions = useMutationOptions(["cluster/list"])
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <Autocomplete
            label={label}
            selected={{key: credId, short: shortUuid(credId)}}
            options={options}
            loading={query.isLoading || updateCluster.isLoading}
            onUpdate={handleUpdate}
        />
    )

    function handleUpdate(option: Option | null) {
        updateCluster.mutate({...cluster, [passKey]: option?.key})
    }

    function handleMemoOptions(): Option[] {
        return Object.entries(query.data ?? {}).map(([key, value]) => ({ key, short: shortUuid(key), name: value.username }))
    }
}

import {CredentialOptions, shortUuid} from "../../../app/utils";
import {useMemo} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {clusterApi, credentialApi} from "../../../app/api";
import {Cluster, Credentials, CredentialType} from "../../../app/types";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {Autocomplete, Option} from "../../view/Autocomplete";

const keys = {
    [CredentialType.POSTGRES]: "postgresId",
    [CredentialType.PATRONI]: "patroniId",
}

type Props = {
    type: CredentialType
    cluster: Cluster
}

export function OverviewSettingsPassword(props: Props) {
    const { type, cluster } = props
    const passKey = keys[type]
    const passId = cluster.credentials[passKey as keyof Credentials] ?? ""
    const { label } = CredentialOptions[type]

    const query = useQuery(["credentials", type], () => credentialApi.list(type))
    const options = useMemo(handleMemoOptions, [query.data])

    const updateMutationOptions = useMutationOptions([["cluster/list"]])
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <Autocomplete
            label={label}
            selected={{key: passId, short: shortUuid(passId)}}
            options={options}
            loading={query.isLoading || updateCluster.isLoading}
            onUpdate={handleUpdate}
        />
    )

    function handleUpdate(option: Option | null) {
        updateCluster.mutate({...cluster, credentials: {[passKey]: option?.key}})
    }

    function handleMemoOptions(): Option[] {
        return Object.entries(query.data ?? {}).map(([key, value]) => ({ key, short: shortUuid(key), name: value.username }))
    }
}

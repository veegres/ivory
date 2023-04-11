import {CredentialOptions, shortUuid} from "../../../app/utils";
import {useMemo} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {clusterApi, passwordApi} from "../../../app/api";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {AutocompleteUuid, Option} from "../../view/././autocomplete/AutocompleteUuid";
import {PasswordType} from "../../../type/password";
import {Cluster, Credentials} from "../../../type/cluster";

const keys = {
    [PasswordType.POSTGRES]: "postgresId",
    [PasswordType.PATRONI]: "patroniId",
}

type Props = {
    type: PasswordType
    cluster: Cluster
}

export function OptionsPassword(props: Props) {
    const { type, cluster } = props
    const passKey = keys[type]
    const passId = cluster.credentials[passKey as keyof Credentials] ?? ""
    const { label } = CredentialOptions[type]

    const query = useQuery(["credentials", type], () => passwordApi.list(type))
    const options = useMemo(handleMemoOptions, [query.data])

    const updateMutationOptions = useMutationOptions([["cluster/list"]])
    const updateCluster = useMutation(clusterApi.update, updateMutationOptions)

    return (
        <AutocompleteUuid
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

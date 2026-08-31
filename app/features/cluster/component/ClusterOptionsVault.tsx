import {useMemo} from "react"

import {AutocompleteUuid, Option} from "../../../shared/component/autocomplete/AutocompleteUuid"
import {getShortUuid,VaultOptions} from "../../../shared/helper/HelperUtils"
import {useRouterVault} from "../../vault/api/VaultHook"
import {VaultType} from "../../vault/api/VaultType"

type Props = {
    type: VaultType,
    selected?: string,
    username?: string,
    onUpdate: (type: VaultType, s?: string) => void,
    error?: boolean,
}

export function ClusterOptionsVault(props: Props) {
    const {type, onUpdate, selected, username, error = false} = props
    const passId = selected ?? ""
    const {label} = VaultOptions[type]

    const query = useRouterVault(type)
    const options = useMemo(handleMemoOptions, [query.data, username])

    return (
        <AutocompleteUuid
            label={label}
            selected={{key: passId, short: getShortUuid(passId)}}
            options={options}
            loading={query.isPending}
            onUpdate={handleUpdate}
            error={error}
        />
    )

    function handleUpdate(option: Option | null) {
        onUpdate(type, option?.key)
    }

    function handleMemoOptions(): Option[] {
        return Object.entries(query.data ?? {})
            .filter(([, value]) => !username || value.username === username)
            .map(([key, value]) => ({
                key,
                short: getShortUuid(key),
                name: value.username
            }))
    }
}

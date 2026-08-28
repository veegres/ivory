import {useMemo} from "react"

import {useRouterVault} from "../../../features/vault/api/VaultHook"
import {VaultType} from "../../../features/vault/api/VaultType"
import {AutocompleteUuid, Option} from "../../../shared/component/autocomplete/AutocompleteUuid"
import {getShortUuid,VaultOptions} from "../../../shared/helper/HelperUtils"

type Props = {
    type: VaultType,
    selected?: string,
    onUpdate: (type: VaultType, s?: string) => void,
    // username locks the selectable vaults to the ones with this username
    // (e.g. a keeper plugin's engine-required database user)
    username?: string,
    error?: boolean,
}

export function OptionsVault(props: Props) {
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
            search={username}
            error={error}
        />
    )

    function handleUpdate(option: Option | null) {
        onUpdate(type, option?.key)
    }

    function handleMemoOptions(): Option[] {
        return Object.entries(query.data ?? {})
            .filter(([, value]) => username === undefined || value.username === username)
            .map(([key, value]) => ({
                key,
                short: getShortUuid(key),
                name: value.username
            }))
    }
}

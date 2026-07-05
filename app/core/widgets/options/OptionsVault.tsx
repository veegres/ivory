import {useMemo} from "react"

import {useRouterVault} from "../../../features/vault/api/VaultHook"
import {VaultType} from "../../../features/vault/api/VaultType"
import {AutocompleteUuid, Option} from "../../../shared/component/autocomplete/AutocompleteUuid"
import {getShortUuid,VaultOptions} from "../../../shared/helper/HelperUtils"

type Props = {
    type: VaultType,
    selected?: string,
    onUpdate: (type: VaultType, s?: string) => void,
}

export function OptionsVault(props: Props) {
    const {type, onUpdate, selected} = props
    const passId = selected ?? ""
    const {label} = VaultOptions[type]

    const query = useRouterVault(type)
    const options = useMemo(handleMemoOptions, [query.data])

    return (
        <AutocompleteUuid
            label={label}
            selected={{key: passId, short: getShortUuid(passId)}}
            options={options}
            loading={query.isPending}
            onUpdate={handleUpdate}
        />
    )

    function handleUpdate(option: Option | null) {
        onUpdate(type, option?.key)
    }

    function handleMemoOptions(): Option[] {
        return Object.entries(query.data ?? {})
            .map(([key, value]) => ({
                key,
                short: getShortUuid(key),
                name: value.username
            }))
    }
}

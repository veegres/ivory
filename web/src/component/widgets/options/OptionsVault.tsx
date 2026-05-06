import {useMemo} from "react"

import {useRouterVault} from "../../../api/vault/hook"
import {VaultType} from "../../../api/vault/type"
import {shortUuid,VaultOptions} from "../../../app/utils"
import {AutocompleteUuid, Option} from "../../view/autocomplete/AutocompleteUuid"

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
            selected={{key: passId, short: shortUuid(passId)}}
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
                short: shortUuid(key),
                name: value.username
            }))
    }
}

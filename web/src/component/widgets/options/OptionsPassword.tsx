import {CredentialOptions, shortUuid} from "../../../app/utils";
import {useMemo} from "react";
import {AutocompleteUuid, Option} from "../../view/autocomplete/AutocompleteUuid";
import {PasswordType} from "../../../api/password/type";
import {useRouterPassword} from "../../../api/password/hook";

type Props = {
    type: PasswordType,
    selected?: string,
    onUpdate: (type: PasswordType, s?: string) => void,
}

export function OptionsPassword(props: Props) {
    const {type, onUpdate, selected} = props
    const passId = selected ?? ""
    const {label} = CredentialOptions[type]

    const query = useRouterPassword(type)
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

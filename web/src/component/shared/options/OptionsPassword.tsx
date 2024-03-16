import {CredentialOptions, shortUuid} from "../../../app/utils";
import {useMemo} from "react";
import {useQuery} from "@tanstack/react-query";
import {PasswordApi} from "../../../app/api";
import {AutocompleteUuid, Option} from "../../view/autocomplete/AutocompleteUuid";
import {PasswordType} from "../../../type/password";

type Props = {
    type: PasswordType,
    selected?: string,
    onUpdate: (type: PasswordType, s?: string) => void,
}

export function OptionsPassword(props: Props) {
    const {type, onUpdate, selected} = props
    const passId = selected ?? ""
    const {label} = CredentialOptions[type]

    const query = useQuery({queryKey: ["credentials", type], queryFn: () => PasswordApi.list(type)})
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

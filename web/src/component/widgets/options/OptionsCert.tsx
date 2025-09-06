import {useMemo} from "react";

import {AutocompleteUuid, Option} from "../../view/autocomplete/AutocompleteUuid";
import {CertOptions, shortUuid} from "../../../app/utils";
import {CertType} from "../../../api/cert/type";
import {useRouterCertList} from "../../../api/cert/hook";

type Props = {
    type: CertType,
    selected?: string,
    onUpdate: (type: CertType, s?: string) => void,
}

export function OptionsCert(props: Props) {
    const {type, selected, onUpdate} = props
    const certId = selected ?? ""
    const {label} = CertOptions[type]

    const query = useRouterCertList(type)
    const options = useMemo(handleMemoOptions, [query.data])

    return (
        <AutocompleteUuid
            label={label}
            selected={{key: certId, short: shortUuid(certId)}}
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
                name: value.fileName
            }))
    }
}

import {useMemo} from "react"

import {useRouterCertList} from "../../../features/cert/api/CertHook"
import {CertType} from "../../../features/cert/api/CertType"
import {AutocompleteUuid, Option} from "../../../shared/component/autocomplete/AutocompleteUuid"
import {CertOptions, getShortUuid} from "../../../shared/helper/HelperUtils"

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
            selected={{key: certId, short: getShortUuid(certId)}}
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
                name: value.fileName
            }))
    }
}

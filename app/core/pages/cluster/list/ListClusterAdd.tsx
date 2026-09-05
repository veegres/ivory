import {Add} from "@mui/icons-material"

import {Feature} from "../../../../features/Feature"
import {ManageAccessBox} from "../../../../features/management/component/ManageAccess"
import {TriggerButton} from "../../../../shared/component/button/TriggerButton"

type Props = {
    onClick: () => void,
    disabled?: boolean,
    withLabel?: boolean,
    size?: number,
}

export function ListClusterAdd(props: Props) {
    const {onClick, disabled = false, withLabel = false, size} = props

    return (
        <ManageAccessBox feature={Feature.ManageClusterCreate}>
            <TriggerButton
                variant={withLabel ? "button_label" : "button"}
                title={"ADD CLUSTER"}
                label={"Add"}
                icon={<Add/>}
                size={size}
                onClick={onClick}
                disabled={disabled}
            />
        </ManageAccessBox>
    )
}

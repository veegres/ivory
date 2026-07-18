import {Add} from "@mui/icons-material"

import {Feature} from "../../../../features/Feature"
import {ManageAccessBox} from "../../../../features/management/component/ManageAccess"
import {TriggerButton} from "../../../../shared/component/button/TriggerButton"

type Props = {
    onClick: () => void,
    disabled?: boolean,
    withLabel?: boolean,
}

export function ListClusterAdd(props: Props) {
    const {onClick, disabled = false, withLabel = false} = props

    return (
        <ManageAccessBox feature={Feature.ManageClusterUpdate}>
            <TriggerButton
                variant={withLabel ? "button_label" : "button"}
                title={"ADD CLUSTER"}
                label={"Add"}
                icon={<Add/>}
                onClick={onClick}
                disabled={disabled}
            />
        </ManageAccessBox>
    )
}

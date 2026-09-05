import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../shared/component/box/InfoColorBoxRow"
import {DateTimeFormatter} from "../../../shared/helper/HelperUtils"
import {UserRegistrationStatus} from "../api/UserType"

const STATUS: {[key in UserRegistrationStatus]: {label: string, color: string}} = {
    [UserRegistrationStatus.ACTIVE]: {label: "registered", color: "success"},
    [UserRegistrationStatus.PENDING]: {label: "link outstanding", color: "info"},
    [UserRegistrationStatus.EXPIRED]: {label: "link expired", color: "warning"},
    [UserRegistrationStatus.MISSING]: {label: "no password", color: "inherit"},
}

type Props = {
    status?: UserRegistrationStatus,
    expiresAt?: string,
}

export function UserBasicStatus(props: Props) {
    const {status, expiresAt} = props

    return (
        <InfoColorBoxRow>{renderTags()}</InfoColorBoxRow>
    )

    function renderTags() {
        if (!status) return <InfoColorBox label={"not configured"} title={"This user doesn't use BASIC auth"}/>
        const {label, color} = STATUS[status]
        return <InfoColorBox label={label} color={color} title={getStatusTitle()}/>
    }

    function getStatusTitle() {
        if (!expiresAt) return
        return `Expires on ${DateTimeFormatter.utc(expiresAt)}`
    }
}
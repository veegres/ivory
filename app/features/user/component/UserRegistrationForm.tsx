import {Stars} from "@mui/icons-material"
import {Box, ToggleButton, Tooltip} from "@mui/material"
import {useState} from "react"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {KeyEnterInput} from "../../../shared/component/input/KeyEnterInput"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {useRouterInfo} from "../../management/api/ManagementHook"
import {useRouterUserCreate} from "../api/UserHook"
import {UserAuthType, UserRegistered, UserRegistration, UserSetupRequest} from "../api/UserType"
import {UserAuthTypes} from "./UserAuthTypes"
import {UserRegistrationLink} from "./UserRegistrationLink"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    row: {display: "flex", alignItems: "center", gap: 1},
    grow: {flexGrow: 1, minWidth: "150px"},
}

type SetupProps = {
    setup: true,
    value: UserSetupRequest,
    onChange: (value: UserSetupRequest) => void,
}

type ManageProps = {
    setup?: false,
    value?: never,
    onChange?: never,
}

type Props = SetupProps | ManageProps

export function UserRegistrationForm(props: Props) {
    const setup = props.setup === true
    const [own, setOwn] = useState<UserSetupRequest>({username: "", password: "", authTypes: []})
    const [superuser, setSuperuser] = useState(false)
    const [registration, setRegistration] = useState<UserRegistration>()

    const info = useRouterInfo()
    const create = useRouterUserCreate(handleCreated)
    const value = props.setup ? props.value : own
    const allowed = setup || info.data?.auth.user?.superuser === true

    return (
        <Box sx={SX.box}>
            <Box sx={SX.row}>
                <Box sx={SX.grow}>
                    <KeyEnterInput
                        label={"Username"}
                        value={value.username}
                        onChange={(e) => handleChange({...value, username: e.target.value})}
                        onEnterPress={handleCreate}
                    />
                </Box>
                {renderSuperuser()}
            </Box>
            <Box sx={SX.row}>
                <UserAuthTypes
                    value={value.authTypes}
                    onChange={(authTypes) => handleChange({...value, authTypes})}
                />
                {setup ? renderPassword() : renderCreate()}
            </Box>
            {!setup && superuser && <AlertCentered severity={"warning"} text={renderSuperuserDescription()}/>}
            {registration && <UserRegistrationLink registration={registration}/>}
        </Box>
    )

    function renderSuperuser() {
        const on = setup || superuser
        return (
            <Tooltip title={getSuperuserTooltip()} placement={"top"} arrow disableInteractive>
                <Box component={"span"}>
                    <ToggleButton
                        value={"superuser"}
                        selected={on}
                        disabled={setup || !allowed}
                        onClick={() => setSuperuser(!superuser)}
                    >
                        <Stars/>
                    </ToggleButton>
                </Box>
            </Tooltip>
        )
    }

    function renderPassword() {
        if (!value.authTypes.includes(UserAuthType.BASIC)) return null
        return (
            <Box sx={SX.grow}>
                <KeyEnterInput
                    label={"Password"}
                    value={value.password}
                    hidden
                    onChange={(e) => handleChange({...value, password: e.target.value})}
                />
            </Box>
        )
    }

    function renderCreate() {
        return (
            <Tooltip title={getCreateTooltip()} placement={"top"} arrow disableInteractive>
                <Box component={"span"} sx={SX.grow}>
                    <SimpleButton
                        fullWidth={true}
                        variant={"contained"}
                        color={"primary"}
                        loading={create.isPending}
                        disabled={!isComplete()}
                        onClick={handleCreate}
                    >
                        Register
                    </SimpleButton>
                </Box>
            </Tooltip>
        )
    }

    function renderSuperuserDescription() {
        return (<>
            A <b>superuser</b> holds <b>every permission</b> and can never lose one, not even to
            themselves - that is what keeps <b>Ivory</b> administrable. Only a superuser can register,
            delete or re-register another superuser. It cannot be changed afterwards: taking it back
            means deleting the user and registering them again.
        </>)
    }

    function getSuperuserTooltip() {
        if (setup) return "The first Ivory user is always a superuser"
        if (!allowed) return "Only a superuser can register another superuser"
        return "Superuser"
    }

    function getCreateTooltip() {
        if (!value.username.trim()) return "A user needs a name"
        if (value.authTypes.length === 0) return "A user needs at least one way to sign in"
        if (value.authTypes.includes(UserAuthType.BASIC)) return "Register this user and issue their registration link"
        return "Register this user"
    }

    function handleChange(next: UserSetupRequest) {
        if (props.setup) props.onChange(next)
        else setOwn(next)
    }

    function handleCreate() {
        if (setup || !isComplete()) return
        create.mutate({username: value.username.trim(), authTypes: value.authTypes, superuser})
    }

    function handleCreated(registered: UserRegistered) {
        setRegistration(registered.registration)
        setOwn({username: "", password: "", authTypes: []})
        setSuperuser(false)
    }

    function isComplete() {
        return !!value.username.trim() && value.authTypes.length > 0
    }
}

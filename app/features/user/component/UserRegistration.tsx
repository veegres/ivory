import {Box} from "@mui/material"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {UserRegistrationForm} from "./UserRegistrationForm"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
}


export function UserRegistration() {
    return (
        <Box sx={SX.box}>
            <AlertCentered text={renderDescription()}/>
            <UserRegistrationForm/>
        </Box>
    )

    function renderDescription() {
        return (
            "Only registered users can sign in with their configured authentication methods. " +
            "Users set their password using the one-time link they receive."
        )
    }
}
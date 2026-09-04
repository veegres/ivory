import {Box, Button} from "@mui/material"
import {ChangeEvent} from "react"

import {useRouterConnect} from "../../../../features/auth/api/AuthHook"
import {OidcConfig} from "../../../../features/config/api/ConfigType"
import {UserAuthType} from "../../../../features/user/api/UserType"
import {KeyEnterInput} from "../../../../shared/component/input/KeyEnterInput"
import {SxPropsMap} from "../../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    config: OidcConfig,
    onChange: (config: OidcConfig) => void,
}

export function ConfigAuthOidc(props: Props) {
    const {config, onChange} = props
    const connect = useRouterConnect()
    return (
        <Box sx={SX.box}>
            <KeyEnterInput
                label={"Redirect URL"}
                value={config.redirectUrl}
                helperText={"This URL must be allowed in your provider configuration"}
                required={false}
                disabled={true}
                onChange={handleConfigChange("redirectUrl")}
            />
            <KeyEnterInput
                label={"Issuer URL"}
                value={config.issuerUrl}
                helperText={"Specify the complete HTTP url, including the protocol prefix http:// or https://"}
                onChange={handleConfigChange("issuerUrl")}
            />
            <KeyEnterInput label={"Client ID"} value={config.clientId} onChange={handleConfigChange("clientId")}/>
            <KeyEnterInput label={"Client secret"} value={config.clientSecret} onChange={handleConfigChange("clientSecret")} hidden/>
            <Button
                color={"success"}
                loading={connect.isPending}
                onClick={() => connect.mutate({type: UserAuthType.OIDC, config})}>
                Test Connection
            </Button>
        </Box>
    )

    function handleConfigChange(key: string) {
        return (e: ChangeEvent<HTMLInputElement>) => {
            onChange({...config, [key]: e.target.value})
        }
    }
}
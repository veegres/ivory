import {Box, Button} from "@mui/material";
import {ChangeEvent} from "react";

import {useRouterConnect} from "../../../../api/auth/hook";
import {AuthType} from "../../../../api/auth/type";
import {LdapConfig} from "../../../../api/config/type";
import {SxPropsMap} from "../../../../app/type";
import {KeyEnterInput} from "../../../view/input/KeyEnterInput";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    config: LdapConfig,
    onChange: (config: LdapConfig) => void,
}

export function ConfigAuthLdap(props: Props) {
    const {config, onChange} = props
    const connect = useRouterConnect()
    return (
        <Box sx={SX.box}>
            <KeyEnterInput
                label={"URL"}
                helperText={"Specify the complete LDAP url, including the protocol prefix ldap:// or ldaps://"}
                value={config.url}
                onChange={handleConfigChange("url")}
            />
            <KeyEnterInput label={"Bind DN"} value={config.bindDN} onChange={handleConfigChange("bindDN")}/>
            <KeyEnterInput label={"Bind password"} value={config.bindPass} onChange={handleConfigChange("bindPass")} hidden/>
            <KeyEnterInput label={"Base DN"} value={config.baseDN} onChange={handleConfigChange("baseDN")}/>
            <KeyEnterInput
                label={"Filter"}
                helperText={"User search filter, for example, '(cn=%s)' or '(sAMAccountName=%s)'"}
                value={config.filter}
                onChange={handleConfigChange("filter")}
            />
            <Button
                color={"success"}
                loading={connect.isPending}
                onClick={() => connect.mutate({type: AuthType.LDAP, config})}>
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
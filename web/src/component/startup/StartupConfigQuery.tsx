import {Switch} from "@mui/material";
import {StartupBlockConfig} from "./StartupBlockConfig";

type Props = {
    onChange: (enable: boolean) => void,
}

export function StartupConfigQuery(props: Props) {
    const {onChange} = props
    return (
        <StartupBlockConfig
            label={"Manual Queries"}
            description={<>
                It will enable manual SQL query creation and system query editing that provided
                by the <b>Ivory</b> as default templates for troubleshooting. <b>Ivory</b> doesn't
                restrict any how SQL query syntax, so, please, keep in mind that it will be allowed to
                create UPDATE, CREATE or DELETE queries as well as SELECT.
            </>}
            severity={"info"}
            renderAction={
                <Switch size={"small"} onChange={e => onChange(e.target.checked)}/>
            }
        />
    )
}

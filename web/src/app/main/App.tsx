import {Box} from "@mui/material";
import {ReactNode} from "react";
import {Header} from "./Header";
import {Footer} from "./Footer";
import {Body} from "./Body";
import {useRouterInfo} from "../../api/management/hook";
import {SxPropsMap} from "../type";
import {AuthType} from "../../api/auth/type";

const SX: SxPropsMap = {
    // NOTE: we need -8px for the bottom scroll, right scroll is always shown that is why we don't need calc
    box: {display: "flex", flexDirection: "column", gap: 2, flexWrap: "nowrap", minHeight: "calc(100vh - 8px)", minWidth: "1010px"},
    body: {display: "flex", justifyContent: "center", alignItems: "stretch", flexGrow: 1},
}

type Props = {
    children?: ReactNode,
}

export function App(props: Props) {
    const info = useRouterInfo()
    const show = !info.isPending && !!info.data && info.data.secret.key && info.data.auth.authorised

    return (
        <Box sx={SX.box}>
            <Box component={"header"}>
                <Header company={info.data?.config.company ?? "VEEGRES"} show={!props.children && show} auth={info.data?.auth.type ?? AuthType.NONE}/>
            </Box>
            <Box component={"main"} sx={SX.body}>
                {props.children ? props.children : <Body info={info}/>}
            </Box>
            <Box component={"footer"}>
                <Footer label={info.data?.version.label ?? "none"} tag={info.data?.version.tag ?? "master"}/>
            </Box>
        </Box>
    );
}

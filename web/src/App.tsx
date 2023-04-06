import {Box} from "@mui/material";
import {Header} from "./component/section/Header";
import {Body} from "./component/section/Body";
import {useQuery} from "@tanstack/react-query";
import {infoApi} from "./app/api";
import {Menu} from "./component/settings/menu/Menu";
import {Footer} from "./component/section/Footer";
import {SxPropsMap} from "./type/common";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2, flexWrap: "nowrap", height: "100vh", minWidth: "1010px"},
    body: {display: "flex", alignItems: "center", justifyContent: "center", flexGrow: 1},
}

export function App() {
    const info = useQuery(["info"], infoApi.get)
    const show = !info.isLoading && !info.isError && info.data.secret.key

    return (
            <Box sx={SX.box}>
                <Box>
                    <Header company={info.data?.company ?? "VEEGRES"} show={show}/>
                </Box>
                <Box sx={SX.body}>
                    <Body info={info}/>
                    <Menu/>
                </Box>
                <Box>
                    <Footer/>
                </Box>
            </Box>
    );
}

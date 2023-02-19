import {Grid} from "@mui/material";
import {Header} from "./component/section/Header";
import {Body} from "./component/section/Body";
import {useQuery} from "@tanstack/react-query";
import {infoApi} from "./app/api";
import {Menu} from "./component/settings/menu/Menu";

export function App() {
    const info = useQuery(["info"], infoApi.get)

    return (
        <Grid sx={{height: "100vh", minWidth: "1010px"}} container direction={"column"} spacing={2} flexWrap={"nowrap"}>
            <Grid item><Header company={info.data?.company ?? "VEEGRES"} show={!info.isLoading && !info.isError}/></Grid>
            <Grid item container flexGrow={1} justifyContent={"center"} alignItems={"center"}>
                <Body info={info}/>
                <Menu/>
            </Grid>
        </Grid>
    );
}

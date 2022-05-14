import React from 'react';
import {Grid} from "@mui/material";
import {Header} from "./component/section/Header";
import {Body} from "./component/section/Body";

export function App() {
    return (
        <Grid sx={{ height: "100vh" }} container direction={"column"} spacing={2} flexWrap={"nowrap"}>
            <Grid item><Header /></Grid>
            <Grid item container flexGrow={1} justifyContent={"center"} alignItems={"center"}>
                <Body />
            </Grid>
        </Grid>
    );
}

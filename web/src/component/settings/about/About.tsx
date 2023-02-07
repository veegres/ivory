import {MenuWrapper} from "../menu/MenuWrapper";
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll";
import {Box, Divider, IconButton} from "@mui/material";
import {SxPropsMap} from "../../../app/types";
import {MenuItemBox} from "../menu/MenuItemBox";
import {MenuItemText} from "../menu/MenuItemText";
import {OpenInNew} from "@mui/icons-material";

const SX: SxPropsMap = {
    scroll: {display: "flex", flexDirection: "column", padding: "0 15px", gap: 3},
    image: {display: "flex", justifyContent: "center"},
    text: {fontSize: 14},
}

export function About() {
    return (
        <MenuWrapper>
            <MenuWrapperScroll sx={SX.scroll}>
                <Box sx={SX.image}><img src={"/postgres.png"} width={200} height={200} alt={"Ivory"}/></Box>
                <Box sx={SX.text}>
                    Ivory is an open-source project which is designed to simplify and visualize work with
                    postgres clusters.
                </Box>
                <MenuItemBox name={"Links"}>
                    <MenuItemText title={"GitHub"} button={renderLink("https://github.com/veegres/ivory")}/>
                    <Divider/>
                    <MenuItemText title={"Docker Hub"} button={renderLink("https://hub.docker.com/r/aelsergeev/ivory")}/>
                    <Divider/>
                    <MenuItemText title={"Contribution & Issues"} button={renderLink("https://github.com/veegres/ivory/issues")}/>
                    <Divider/>
                    <MenuItemText title={"Releases"} button={renderLink("https://github.com/veegres/ivory/releases")}/>
                </MenuItemBox>
            </MenuWrapperScroll>
        </MenuWrapper>
    )

    function renderLink(href: string) {
        return (
            <IconButton href={href} target={"_blank"} rel={"noopener"}>
                <OpenInNew/>
            </IconButton>
        )
    }
}

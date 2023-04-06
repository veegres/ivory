import {MenuWrapper} from "../menu/MenuWrapper";
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll";
import {Box, IconButton} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {MenuItemBox} from "../menu/MenuItemBox";
import {MenuItemText} from "../menu/MenuItemText";
import {OpenInNew} from "@mui/icons-material";
import {IvoryLinks} from "../../../app/utils";

const SX: SxPropsMap = {
    scroll: {display: "flex", flexDirection: "column", padding: "0 15px", gap: 3},
    image: {display: "flex", justifyContent: "center"},
    text: {fontSize: 14},
}

export function About() {
    return (
        <MenuWrapper>
            <MenuWrapperScroll sx={SX.scroll}>
                <Box sx={SX.image}><img src={"/ivory.png"} width={200} height={200} alt={"Ivory"}/></Box>
                <Box sx={SX.text}>
                    Ivory is an open-source project which is designed to simplify and visualize work with
                    postgres clusters.
                </Box>
                <MenuItemBox name={"Links"}>
                    {Object.values(IvoryLinks).map(({name, link}) => (
                        <MenuItemText title={name} button={renderLink(link)}/>
                    ))}
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

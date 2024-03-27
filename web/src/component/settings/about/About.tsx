import {MenuWrapper} from "../menu/MenuWrapper";
import {Box, IconButton} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {OpenInNew} from "@mui/icons-material";
import {IvoryLinks} from "../../../app/utils";
import {List} from "../../view/box/List";
import {ListItem} from "../../view/box/ListItem";

const SX: SxPropsMap = {
    scroll: {display: "flex", flexDirection: "column", padding: "0 15px", gap: 3},
    image: {display: "flex", justifyContent: "center"},
    text: {fontSize: 14, textAlign: "justify"},
}

export function About() {
    return (
        <MenuWrapper>
            <Box sx={SX.scroll}>
                <Box sx={SX.image}><img src={"/ivory.png"} width={200} height={200} alt={"Ivory"}/></Box>
                <Box sx={SX.text}>
                    Ivory is an open-source project designed to simplify and visualize work with Postgres
                    clusters. Initially, this tool was developed to ease the life of developers who maintain
                    Postgres. But I hope it will help manage and troubleshoot Postgres clusters for both
                    developers and database administrators.
                </Box>
                <List name={"Links"}>
                    {Object.values(IvoryLinks).map(({name, link}) => (
                        <ListItem key={name} title={name} button={renderLink(link)}/>
                    ))}
                </List>
            </Box>
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

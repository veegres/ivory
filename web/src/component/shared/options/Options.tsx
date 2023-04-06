import {OptionsPassword} from "./OptionsPassword";
import {PasswordType} from "../../../type/password";
import {Divider} from "@mui/material";
import {OptionsCert} from "./OptionsCert";
import {CertType} from "../../../type/cert";
import {OptionsTags} from "./OptionsTags";
import {Cluster} from "../../../type/cluster";

type Props = {
    cluster: Cluster
}

export function Options(props: Props) {
    const {cluster} = props
    return (
        <>
            <OptionsPassword type={PasswordType.POSTGRES} cluster={cluster}/>
            <OptionsPassword type={PasswordType.PATRONI} cluster={cluster}/>
            <Divider variant={"middle"}/>
            <OptionsCert type={CertType.CLIENT_CA} cluster={cluster}/>
            <OptionsCert type={CertType.CLIENT_CERT} cluster={cluster}/>
            <OptionsCert type={CertType.CLIENT_KEY} cluster={cluster}/>
            <Divider variant={"middle"}/>
            <OptionsTags cluster={cluster}/>
        </>
    )
}

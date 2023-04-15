import {OptionsPassword} from "./OptionsPassword";
import {PasswordType} from "../../../type/password";
import {Divider} from "@mui/material";
import {OptionsCert} from "./OptionsCert";
import {CertType} from "../../../type/cert";
import {OptionsTags} from "./OptionsTags";
import {ClusterOptions} from "../../../type/cluster";
import {CertOptions, CredentialOptions} from "../../../app/utils";

type Props = {
    cluster: ClusterOptions,
    onUpdate: (cluster: ClusterOptions) => void,
}

export function Options(props: Props) {
    const {onUpdate, cluster} = props
    const {credentials, tags, certs} = cluster

    return (
        <>
            <OptionsPassword type={PasswordType.POSTGRES} selected={credentials.postgresId} onUpdate={handlePasswordUpdate}/>
            <OptionsPassword type={PasswordType.PATRONI} selected={credentials.patroniId} onUpdate={handlePasswordUpdate}/>
            <Divider variant={"middle"}/>
            <OptionsCert type={CertType.CLIENT_CA} selected={certs.clientCAId} onUpdate={handleCertUpdate}/>
            <OptionsCert type={CertType.CLIENT_CERT} selected={certs.clientCertId} onUpdate={handleCertUpdate}/>
            <OptionsCert type={CertType.CLIENT_KEY} selected={certs.clientKeyId} onUpdate={handleCertUpdate}/>
            <Divider variant={"middle"}/>
            <OptionsTags selected={tags} onUpdate={handleTagsUpdate}/>
        </>
    )

    function handlePasswordUpdate(t: PasswordType, s?: string) {
        onUpdate({...cluster, credentials: { [CredentialOptions[t].key]: s }})
    }

    function handleCertUpdate(t: CertType, s?: string) {
        onUpdate({...cluster, certs: { [CertOptions[t].key]: s }})
    }

    function handleTagsUpdate(tags: string[]) {
        onUpdate({...cluster, tags})
    }
}

import {OptionsPassword} from "./OptionsPassword";
import {PasswordType} from "../../../api/password/type";
import {Divider, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {OptionsCert} from "./OptionsCert";
import {CertType} from "../../../api/cert/type";
import {OptionsTags} from "./OptionsTags";
import {ClusterOptions} from "../../../api/cluster/type";
import {CertOptions, CredentialOptions} from "../../../app/utils";

import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    tls: {width: "50%", fontWeight: "bold"},
}

type Props = {
    cluster: ClusterOptions,
    onUpdate: (cluster: ClusterOptions) => void,
}

export function Options(props: Props) {
    const {onUpdate, cluster} = props
    const {credentials, tags, certs, tls} = cluster

    return (
        <>
            <OptionsPassword type={PasswordType.POSTGRES} selected={credentials.postgresId} onUpdate={handlePasswordUpdate}/>
            <OptionsPassword type={PasswordType.PATRONI} selected={credentials.patroniId} onUpdate={handlePasswordUpdate}/>
            <Divider variant={"middle"}/>
            <ToggleButtonGroup size={"small"} fullWidth>
                <ToggleButton onClick={handleTlsSidecarUpdate} selected={tls.sidecar} value={"sidecar"}>Sidecar</ToggleButton>
                <ToggleButton sx={SX.tls} disabled={true} value={"tls"}>TLS</ToggleButton>
                <ToggleButton onClick={handleTlsDatabaseUpdate} selected={tls.database} value={"database"}>Database</ToggleButton>
            </ToggleButtonGroup>
            <Divider variant={"middle"}/>
            <OptionsCert type={CertType.CLIENT_CA} selected={certs.clientCAId} onUpdate={handleCertUpdate}/>
            <OptionsCert type={CertType.CLIENT_CERT} selected={certs.clientCertId} onUpdate={handleCertUpdate}/>
            <OptionsCert type={CertType.CLIENT_KEY} selected={certs.clientKeyId} onUpdate={handleCertUpdate}/>
            <Divider variant={"middle"}/>
            <OptionsTags selected={tags} onUpdate={handleTagsUpdate}/>
        </>
    )

    function handlePasswordUpdate(t: PasswordType, s?: string) {
        onUpdate({...cluster, credentials: {...cluster.credentials, [CredentialOptions[t].key]: s}})
    }

    function handleCertUpdate(t: CertType, s?: string) {
        onUpdate({...cluster, certs: {...cluster.certs, [CertOptions[t].key]: s}})
    }

    function handleTagsUpdate(tags: string[]) {
        onUpdate({...cluster, tags})
    }

    function handleTlsSidecarUpdate() {
        onUpdate({...cluster, tls: {...cluster.tls, sidecar: !tls.sidecar}})
    }

    function handleTlsDatabaseUpdate() {
        onUpdate({...cluster, tls: {...cluster.tls, database: !tls.database}})
    }
}

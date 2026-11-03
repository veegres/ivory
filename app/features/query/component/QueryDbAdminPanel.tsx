import { Alert, Box } from "@mui/material";
import { ReactNode } from "react";

import { Refresher } from "../../../core/widgets/browser/Refresher";
import { TitleBox } from "../../../shared/component/box/TitleBox";
import { SxPropsMap } from "../../../shared/helper/HelperType";
import { Feature } from "../../Feature";
import {
  ManageAccess,
  ManageAccessBox,
} from "../../management/component/ManageAccess";
import { QueryApi } from "../api/QueryRouter";
import { Connection } from "../api/QueryType";

export const QueryDbAdminSX: SxPropsMap = {
  panel: {
    display: "flex",
    flexDirection: "column",
    gap: 1,
    minHeight: 0,
    height: "100%",
  },
  toolbar: {
    display: "flex",
    alignItems: "center",
    justifyContent: "flex-end",
  },
  list: {
    display: "flex",
    flexDirection: "column",
    gap: 1,
    overflow: "auto",
    minHeight: 0,
    flexGrow: 1,
    paddingRight: 0.5,
  },
  cardHead: {
    display: "flex",
    flexWrap: "wrap",
    alignItems: "center",
    justifyContent: "space-between",
    gap: 1,
  },
  cardActions: {
    display: "flex",
    alignItems: "center",
    gap: 0.5,
    flexShrink: 0,
  },
  chips: { display: "flex", flexWrap: "wrap", gap: 0.5, alignItems: "center" },
  metaGrid: {
    display: "grid",
    gridTemplateColumns: { xs: "1fr", sm: "repeat(3, 1fr)" },
    gap: 1,
    fontSize: "13px",
    color: "text.secondary",
  },
  metaLabel: {
    fontSize: "11px",
    textTransform: "uppercase",
    color: "text.disabled",
    lineHeight: 1.2,
  },
  metaValue: { fontFamily: "monospace", color: "text.primary" },
  form: { paddingTop: 0.5 },
};

type Props = {
  connection: Connection;
  viewFeature: Feature;
  createFeature: Feature;
  alert: string;
  createLabel: string;
  listKind: "role" | "database";
  renderCreate: ReactNode;
  children: ReactNode;
};

export function QueryDbAdminPanel(props: Props) {
  const {
    connection,
    viewFeature,
    createFeature,
    alert,
    createLabel,
    listKind,
    renderCreate,
    children,
  } = props;
  const queryKey =
    listKind === "role"
      ? QueryApi.dbRoleList.key(connection)
      : QueryApi.dbDatabaseAdminList.key(connection);

  return (
    <ManageAccessBox
      sx={QueryDbAdminSX.panel}
      feature={viewFeature}
      error={true}
    >
      <Alert severity={"info"} sx={{ flexShrink: 0 }}>
        {alert}
      </Alert>
      <Box sx={QueryDbAdminSX.toolbar}>
        <Refresher size={26} queryKeys={[queryKey]} />
      </Box>
      <ManageAccess feature={createFeature}>
        <TitleBox
          label={createLabel}
          island={true}
          dense={true}
          collapsible={true}
          defaultOpen={true}
        >
          <Box sx={QueryDbAdminSX.form}>{renderCreate}</Box>
        </TitleBox>
      </ManageAccess>
      <Box sx={QueryDbAdminSX.list}>{children}</Box>
    </ManageAccessBox>
  );
}

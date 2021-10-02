import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {RootState} from '../../app/store';
import {restClient} from "../../app/hooks";

export interface ClusterList {
    name: string,
    nodes: string[]
}

export interface ClusterListState {
    data: ClusterList[];
    error?: string;
    status: 'idle' | 'loading' | 'failed';
}

const initialState: ClusterListState = {
    data: [],
    status: 'idle',
};

export const incrementAsync = createAsyncThunk(
    'cluster/fetchList',
    async (data, { rejectWithValue }) => {
        try {
            return await restClient<ClusterList[]>(`/cluster/list`)
        } catch (error) {
            return rejectWithValue(error.response.data)
        }
    }
);

export const clusterListSlice = createSlice({
    name: 'cluster/list',
    initialState,
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(incrementAsync.pending, (state) => {
                state.status = 'loading';
            })
            .addCase(incrementAsync.fulfilled, (state, action) => {
                state.status = 'idle';
                state.data = action.payload.response;
            })
            .addCase(incrementAsync.rejected, (state, action) => {
                state.status = 'idle';
                state.error = action.error.message;
            });
    },
});

export const selectClusterListData = (state: RootState) => state.clusterList.data;

import React  from "react";
import Layout from "../../lib/components/layout";
import {auth} from "../../lib/api";
import {Error} from "../../lib/components/controls"
import {useRouter} from "../../lib/hooks/router";
import {useAPI} from "../../lib/hooks/api";

const LoginTokenPage = () => {
    const router = useRouter();
    const { response, error, loading } = useAPI(() => {
        if (!router.query.t) {
            throw new Error("Missing login token!");
        }
        return auth.login_token(router.query.t);
    });
    if (loading) {
        return null;
    }
    if (!error && response) {
        router.push(router.query.next ? router.query.next : '/');
    }

    return (
        <Layout logged={false}>
        {!!error && <Error error={error}/>}
        </Layout>
    );
};

export default LoginTokenPage;

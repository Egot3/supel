import { fail, redirect } from "@sveltejs/kit";
import axios from "axios";

/**
 @satisfies {import{'./$types'}.actions;}
 */
export const actions = {
    register: async ({request,url}) => {
        const data = await request.formData()

        const email = data.get("email")
        const password = data.get("password")
        const nickname = data.get("nickname")

        if (!email){
            return fail(422, {email, missing: true})
        }
        if (!password){
            return fail(422, {password, missing: true})
        }
        const passwordString = password?.toString()


        let matchResult = /(\d+)/.exec(passwordString)
        if ( matchResult===null || !(matchResult[1].length > 4)) {
            return fail(422, {password, digitRequired: true})
        }

        matchResult = /\s+/.exec(passwordString)
        if(matchResult!==null){
            return fail(422, {password, whitespace: true})
        }

        matchResult = /([A-Z])+/.exec(passwordString)
        if ( matchResult===null || !(matchResult[1].length > 1)) {
            return fail(422, {password, upperCaseRequired: true})
        }

        matchResult = /([a-z])+/.exec(passwordString)
        if ( matchResult===null || !(matchResult[1].length > 1)) {
            return fail(422, {password, lowerCaseRequired: true})
        }

        const tokenResponse = await axios.post("http://localhost:5003/registration", {
            email: email,
            password: password,
            nickname: nickname==null?email:nickname,
        })
        if(tokenResponse.data == null){
            return fail(500)
        }
        if(tokenResponse.status===500){
            return fail(500, {tokenResponse})
        }

        const token = JSON.parse(tokenResponse.data)["token"]

        localStorage.setItem("token", token)

        if (url.searchParams.has('redirectTo')){
            redirect(303, url.searchParams.get('redirectTo')!)
        }

        return {success:true}
    }
}
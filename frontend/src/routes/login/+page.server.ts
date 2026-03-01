import { fail, redirect, type Actions } from "@sveltejs/kit";
import axios from "axios";

/**
 @satisfies {import{'./$types'}.actions;}
 */
export const actions: Actions = {
    login: async ({request, url}) => {//ох, зря я сюда полез

        console.log("Login fired")

        const data = await request.formData()
        const email = data.get("email")
        const password = data.get("password")//перехват пароля уже проблема пользователя

        if (!email){
            return fail(422, {email, missing: true})
        }
        if (!password){
            return fail(422, {password, missing: true})
        }
        const passwordString = password?.toString()
        console.log(passwordString)


        const digitCount = (passwordString.match(/\d/g) || []).length;
        if (digitCount < 4) {
            return fail(422, {password, digitRequired: true});
        }
        console.log("passed digit check")

        let matchResult = /\s+/.exec(passwordString)
        if(matchResult!==null){
            return fail(422, {password, whitespace: true})
        }
        console.log("passed space check")

        matchResult = /([A-Z])+/.exec(passwordString)
        if ( matchResult===null || !(matchResult[1].length >= 1)) {
            return fail(422, {password, upperCaseRequired: true})
        }
        console.log("passed Uppercase check")


        matchResult = /([a-z])+/.exec(passwordString)
        if ( matchResult===null || !(matchResult[1].length >= 1)) {
            return fail(422, {password, lowerCaseRequired: true})
        }
        console.log("passed lowercase check")

        
        const tokenResponse = await axios.post("http://localhost:5003/login", {
            email: email,
            password: password,
        })
        if(tokenResponse.data == null){
            return fail(500)
        }
        if(tokenResponse.status==401){
            return fail(401, {tokenResponse})
        }
        const token = tokenResponse.data.token

        if (url.searchParams.has('redirectTo')){
            redirect(303, url.searchParams.get('redirectTo')!)
        }

        //localStorage.setItem("token", token)

        return {
            success: true,
            token: token,
        }
    },
}
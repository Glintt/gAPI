
const OAUTH_API = require("@/api/auth");
import Vue from "vue";

export const vmA = new Vue({
    data: {
        loggedIn : false,
        user: null
    },
    methods:{
        logout(){
            OAUTH_API.clearAccessToken();
            this.loggedIn = false
        },
        isLoggedIn(){
            this.loggedIn = OAUTH_API.getToken() && new Date().getTime() < OAUTH_API.getExpirationTime();
            return this.loggedIn;
        },
        authenticate(user, cb){
            OAUTH_API.authenticate(user.username,user.password, (response) => {
                cb(response)
            });
        },
        currentUser() {
            return this.user
        }
    }
})
/* 
export function logout() {
    OAUTH_API.clearAccessToken();    
}
  
export function isLoggedIn() {
    return OAUTH_API.getToken() && new Date().getTime() < OAUTH_API.getExpirationTime();
}
*/

export function updateAuthenticatedStatus(){

}
export function requireAuth(to, from, next) {
    if (!vmA.isLoggedIn()) {
        next({
            path: '/login'
          });
    } else {
      next();
    }
} 
  
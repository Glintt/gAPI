<template>
    <div class="container">
        <div class="row">
            <div class="col-md-12">
                <div class="row">
                    <div class="col-md-6 mx-auto">
                        <div class="card rounded-0">
                            <div class="card-header">
                                <h3 class="mb-0">Login</h3>
                            </div>
                            <div class="card-body">
                                <information-panel v-if="loginMessage.showing" :msg="loginMessage.alertText" :className="loginMessage.alertClass"></information-panel>

                                <form class="form" v-on:keyup.13="authenticate">
                                    <div class="form-group">
                                        <label for="uname1">Username</label>
                                        <input type="text" 
                                            v-model="user.username"
                                            class="form-control form-control-lg rounded-0" 
                                            name="uname1" 
                                            id="uname1"
                                            required="">
                                    </div>
    
                                    <div class="form-group">
                                        <label>Password</label>
                                        <input type="password" 
                                            v-model="user.password"
                                            class="form-control form-control-lg rounded-0" 
                                            id="pwd1" required=""
                                            autocomplete="password">
                                    </div>
    
                                    <button type="button" @click="authenticate" class="btn btn-success btn-lg float-right" id="btnLogin">Login</button>
                                </form>    
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import InformationPanel from "@/components/InformationPanel"
    var OAuthAPI = require("@/api/auth")
    export default {
        name: "login",
        data(){
            return {
                user:{
                    username:"",
                    password:""
                },
                loginMessage:{
                    alertClass: "alert alert-success",
                    alertText:"",
                    showing: false
                }
            }
        },
        methods:{
            authenticate : function() {
                this.$oauthUtils.vmA.authenticate(this.user, (response) => {
                    this.loginMessage.showing = true;

                    if(response.status != 200) {
                        this.loginMessage.alertText = "Invalid login.";
                        this.loginMessage.alertClass = "alert alert-danger";
                        return ;
                    }
                    this.loginMessage.alertText = "Login successful";
                    this.loginMessage.alertClass = "alert alert-success";

                    setTimeout(() => {
                        window.location = "/";
                    }, 500);
                });
            }
        },
        components: {  
            InformationPanel    
        }    
    };
</script>

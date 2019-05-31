import { Injectable } from '@angular/core';
import {LoginRequest, LogoutRequest, RegisterRequest, VerifyRequest} from './user.model';
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class UsersApiService {
  registerFormData: RegisterRequest;
  verifyFormData: VerifyRequest;
  loginFormData: LoginRequest;
  logoutFormData: LogoutRequest;
  readonly rootURL ="http://localhost";
  constructor(private http : HttpClient) { }


  postRegister() {
    return this.http.post(this.rootURL+'/register',this.registerFormData)
  }

  postVerify() {
    return this.http.post(this.rootURL+'/verify',this.verifyFormData)
  }

  postLogin() {
    return this.http.post(this.rootURL+'/login',this.loginFormData)
  }

  postLogout() {
    return this.http.post(this.rootURL+'/logout',this.logoutFormData)
  }
}

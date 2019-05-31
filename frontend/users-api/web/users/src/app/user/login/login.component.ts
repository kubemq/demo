import { Component, OnInit } from '@angular/core';
import { UsersApiService } from '../../shared/users-api.service';
import {NgForm} from "@angular/forms";
import {ApiResponse, User} from "../../shared/user.model";
import { ToastrService } from 'ngx-toastr';
@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  constructor(private apiService: UsersApiService, private toastr: ToastrService) {}

  ngOnInit() {
    this.resetForm()
  }
  onSubmit(form: NgForm) {

    this.apiService.postLogin().subscribe((response:ApiResponse)=> {
      if (response.is_error) {
        this.toastr.error(response.message,'Login');
        return
      } else {
        this.toastr.success("user logged in successfully",'Login');
      }
    })
  }
  resetForm(form?: NgForm) {
    if (form != null)
      form.resetForm();

    this.apiService.loginFormData = {
      name: '',
      password:'',
    }
  }

}

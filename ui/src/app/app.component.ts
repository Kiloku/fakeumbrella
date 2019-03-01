import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';

import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit{

	customers : Observable<any[]>;
	constructor(private http:HttpClient){}
	ngOnInit(){
		this.customers = this.http.get<any[]>("http://localhost:8080/customers/")
	}


}
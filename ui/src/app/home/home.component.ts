import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Observable} from "rxjs/Observable";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

	forecasts$ : Observable<any[]>;
	forecasts : object[];
	constructor(private http:HttpClient){}

	ngOnInit(){
		this.getForecasts();
		
	}
	getForecasts(){
        this.http.get("http://localhost:8080/forecasts/customers/").subscribe((res : any)=>{
            //console.log(res);
            this.forecasts = res;
            console.log(this.forecasts);
        });
	}

	parseTime(time){
		//let d : Date;
		let d = new Date(0);
		d.setUTCSeconds(time);

		return d;
	}

  	

}

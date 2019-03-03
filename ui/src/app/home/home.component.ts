import { Component, OnInit, ElementRef, AfterViewInit, ViewChild } from '@angular/core';
import * as BABYLON from 'babylonjs';
import { HttpClient } from '@angular/common/http';
import {Observable} from "rxjs/Observable"; 
import { Chart } from '../chart';
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
	//forecasts$ : Observable<any[]>;
	forecasts : object[];
	customers : object[];
	chart : Chart;
	constructor(private http:HttpClient){}

	ngOnInit(){
		this.getForecasts();
		this.http.get("http://localhost:8080/customers/").subscribe((res : any) => {
			this.customers = res;
			console.log(this.customers);
		});
	}

	ngAfterViewInit() {
		this.chart = new Chart('renderCanvas');
		this.chart.createScene();
	}

	getForecasts(){
        this.http.get("http://localhost:8080/forecasts/customers/").subscribe((res : any)=>{
            this.forecasts = res;
            console.log(this.forecasts);
			this.chart.render(this.customers, this.forecasts);
        });

	}

	parseTime(time){
		let d = new Date(0);
		d.setUTCSeconds(time);

		return d;
	}


}

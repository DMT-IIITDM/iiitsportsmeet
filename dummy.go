package main
  
	import (
		"net/http"
		"github.com/julienschmidt/httprouter"
		"fmt"
		"github.com/ttacon/chalk"
		"gopkg.in/mgo.v2/bson"
		"gopkg.in/mgo.v2"
		"log"
		"html/template"
		"time"
		"strconv"
		"os"
	)
 

	type Idea struct{
		ID     bson.ObjectId `bson:"_id,omitempty" json:"-"`
		Body string	`json:"body" bson:"body"`
		Author string `json:"author" bson:"author"`
		Upvotes int  `json:"upvotes" bson:"upvotes"`
		Iplist []string	`json:"iplist" bson:"iplist"`
		Postip string `json:"postip" bson:"postip"`
		TimeAdded string `json:"timeadded" bson:"timeadded"`
	}

	type BigOb struct{

		IdeaObj []Idea `json:"ideaobj" bson:"ideaobj"`
		Today string `json:"today" bson:"today"`

	}

	type Institute struct{
		ID     bson.ObjectId `bson:"_id,omitempty" json:"-"`
		Iname string `json:"iname" bson:"iname"`
		Uname string `json:"uname" bson:"uname"`
		Phno string `json:"phno" bson:"phno"`

		Password string `json:"password" bson:"password"`
		TimeAdded string `json:"timeadded" bson:"timeadded"`

	}

	type Events struct{
		ID     bson.ObjectId `bson:"_id,omitempty" json:"-"`
		Iname string `json:"iname" bson:"iname"`
		Uname string `json:"uname" bson:"uname"`

		Eventsmale []string `json:"events-male" bson:"event-male"`

		Eventsfemale []string `json:"events-female" bson:"event-female"`

		Males int `json:"male" bson:"male"`

		Females int `json:"female" bson:"female"`
	}
	func ServeHTMl(w http.ResponseWriter, r *http.Request, _ httprouter.Params){


			ip := r.RemoteAddr

			file:= r.URL.Path

			filename:=file[1:]

			if len(filename)==0 {
				Home(w,r,nil)
				fmt.Println(chalk.Yellow,ip," requested home page...",chalk.Reset)
			} else if filename == "register" {
				
				message := ""
				Register(w,r,nil,message)


			} else{  
				http.ServeFile(w,r,filename+".html")

				fmt.Println(chalk.Yellow,ip," requested",filename,"page...",chalk.Reset)	
			}
		
	}


	func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params,message string){

		t,err := template.ParseFiles("register.html")
		riperr(err)
		t.Execute(w,message)

	}

	func Review(w http.ResponseWriter, r *http.Request, _ httprouter.Params, e Events){

		t,err:= template.ParseFiles("review.html")
		riperr(err)
		t.Execute(w,e)

	}


	func Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params){

			
			iname := r.FormValue("iname")
			uname := r.FormValue("uname")
			phno := r.FormValue("phno")

			password := r.FormValue("password")

			year,_,day := time.Now().Date()

			month := time.Now().Month().String()


			t := strconv.Itoa(day)+" "+month+" "+strconv.Itoa(year) 

			lookFor:= Institute{}

		i1 := &Institute{
			Iname : iname,
			Uname : uname,
			Phno : phno,
			Password : password,
			TimeAdded : t,
			}		

		output:= "Your idea is recorded "+i1.Uname

		fmt.Println(output,t)

		if i1.Uname == "" || i1.Iname=="" || i1.Phno=="" || i1.Password=="" {

        	Register(w,r,nil,"Please fill in all details")

        	return
        }


         session, err := mgo.Dial("mongodb://localhost:27017")

        riperr(err)


       	f := session.DB("sports").C("institutes")


		FindWith := bson.M{"uname":uname}

        err = f.Find(FindWith).One(&lookFor)

        fmt.Println("hello"+lookFor.Uname+"hello")

        

       	if lookFor.Uname == ""{
       		c := session.DB("sports").C("institutes")

        err = c.Insert(i1)

        riperr(err)	

        Register(w,r,nil,"You have been Successfully Registered, Please Login")
        return
       	} else{
       		Register(w,r,nil,"This user already exists")	

       	}

       	        defer session.Close()


	}


	func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params){

			
			uname := r.FormValue("uname")
			password := r.FormValue("password")

			year,_,day := time.Now().Date()

			month := time.Now().Month().String()


			t := strconv.Itoa(day)+" "+month+" "+strconv.Itoa(year) 		

		output:= "Your name is recorded "+uname

		fmt.Println(output,t)

		if uname == "" || password=="" {

        	Register(w,r,nil,"Please fill in all details")

        	return
        }
		
		lookFor:= Institute{}

		
        session, err := mgo.Dial("mongodb://localhost:27017")

        riperr(err)

        c := session.DB("sports").C("institutes")

        FindWith := bson.M{"uname":uname}

        err = c.Find(FindWith).One(&lookFor)


        if lookFor.Uname == "" {	

        Register(w,r,nil,"Wrong user id or password")
        return
       	}
        if( lookFor.Password == password){
        	fmt.Println(lookFor)

        	s := session.DB("sports").C("events")

        	FindWith := bson.M{"uname":uname}

        	another_look := Events{}

        	err = s.Find(FindWith).One(&another_look)


        	if(another_look.Uname ==""){
        		Dashboard(w,r,p,lookFor)
        		return
        		
        	}else{

        		Review(w,r,nil,another_look)
        		return
               	}
        	
        } else{
        Register(w,r,nil,"Wrong user id or password")
    	}

     session.Close()

	}


	func Dashboard(w http.ResponseWriter, r *http.Request, p httprouter.Params,i Institute){

		t,err:= template.ParseFiles("dashboard.html")

		riperr(err)

		t.Execute(w,i)

	}


	func ServeEvent(w http.ResponseWriter, r *http.Request, _ httprouter.Params){


			ip := r.RemoteAddr

			filepath:= r.URL.Path

			filename:=filepath[8:]
			if filename == "dmsports.go"{
				http.Redirect(w,r,"/",301)
				fmt.Println(chalk.Yellow,ip," tried to access golang code",chalk.Reset)
				}
			if len(filename)==0 {
				Home(w,r,nil)
				fmt.Println(chalk.Yellow,ip," requested home page...",chalk.Reset)
			} else{  
				http.ServeFile(w,r,filename)

				fmt.Println(filename)	
			}
		
	}
	func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

		

					t,_ := template.ParseFiles("iiitsport.html")
		
					year,_,day := time.Now().Date()

					month := time.Now().Month().String()


					today:= strconv.Itoa(day)+" "+month+" "+strconv.Itoa(year)
					
					bigob1:= BigOb{
					Today : today,
					}
					t.Execute(w,bigob1)
					fmt.Println(time.Now())

	}

	func GoBack(w http.ResponseWriter, r *http.Request, p httprouter.Params){

		http.Redirect(w, r, "/", 301)
	}



	 func Submit(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	 




	 		males,err:= strconv.Atoi(r.FormValue("males"))

	 			riperr(err)

	 		females,err := strconv.Atoi(r.FormValue("females"))

	 			riperr(err)

			r.ParseForm()
			eventmale := r.Form["event-male"]
			eventfemale := r.Form["event-female"]

			name := p.ByName("name")


			e := &Events{

				Uname : name,
				Eventsmale : eventmale,
				Eventsfemale : eventfemale,
				Males : males,
				Females : females,
			}


        session, err := mgo.Dial("mongodb://localhost:27017")
        riperr(err)

        c := session.DB("sports").C("events")

        FindWith := bson.M{"uname":name}

        another_look := Events{}

        err = c.Find(FindWith).One(&another_look)


        	if(another_look.Uname ==""){
        		err = c.Insert(e)

        		riperr(err)

       			 t,err := template.ParseFiles("review.html")

        		riperr(err)

				t.Execute(w,e)

        		
        	}else{
        		
        		UpdateWith := bson.M{"uname": name}
				change := bson.M{"$set": bson.M{"event-male" : eventmale,
				"event-female" : eventfemale,
				"male" : males,
				"female" : females} }
				err = c.Update(UpdateWith, change)

				riperr(err)

				 t,err := template.ParseFiles("review.html")

        		riperr(err)

				t.Execute(w,e)
               	}
        	
        
        
	session.Close()

			
	} 

func Update(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	 




	 		password := r.FormValue("password")


			uname := p.ByName("uname")


		if password=="" {

        	Register(w,r,nil,"Please fill in all details")
        	return
        }

        session, err := mgo.Dial("mongodb://localhost:27017")
        riperr(err)

        f := session.DB("sports").C("institutes")


       lookFor:= Institute{}

		FindWith := bson.M{"uname":uname}
        
       err = f.Find(FindWith).One(&lookFor)


		if ( lookFor.Password == ""){	

        	Register(w,r,nil,"Wrong password")
        	session.Close()
        	return
       	}else if( lookFor.Password == password){
        	fmt.Println(lookFor)

        	Dashboard(w,r,p,lookFor)

        }

   		 session.Close()

	} 
	func main(){
		
		newfeed:= chalk.Red.NewStyle().WithBackground(chalk.Black)

		Server:= httprouter.New()

		

		Server.POST("/submit/:name",Submit)

		Server.POST("/signup",Signup)

		Server.POST("/login/",Login)

		Server.GET("/submit/",GoBack)

		Server.POST("/update/:uname",Update);

        Server.ServeFiles("/resources/*filepath", http.Dir("./resources"))
        Server.ServeFiles("/login/resources/*filepath", http.Dir("./resources"))
        
        Server.GET("/register",ServeHTMl)

        Server.GET("/contacts",ServeHTMl)

        Server.GET("/rules",ServeHTMl)

        Server.GET("/",ServeHTMl)

        Server.GET("/events/*filepath",ServeEvent)

		fmt.Println(newfeed,"waiting at :4747",chalk.Reset);

		http.ListenAndServe(GetPort(),Server)
	}
 
func GetPort() string {
        var port = os.Getenv("PORT")
        // Set a default port if there is nothing in the environment
        if port == "" {
                port = "4747"
                fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
        }
        return ":" + port
}

 func riperr(err error){
 	if err!= nil{

    		log.Fatal(err)

    	}
 }

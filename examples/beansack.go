package examples

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/soumitsalman/coffeemaker/sdk/beansack"
	"github.com/soumitsalman/coffeemaker/sdk/beansack/nlp"
	datautils "github.com/soumitsalman/data-utils"
)

func Retrieval() {
	// initialize the services
	if err := beansack.InitializeBeanSack(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("EMBEDDER_URL"), getEmbedderCtx(), os.Getenv("LLMSERVICE_API_KEY")); err != nil {
		log.Fatalln("initialization not working", err)
	}

	options := beansack.NewSearchOptions().WithTimeWindow(30).WithURLs([]string{
		"https://hackaday.com/2024/05/13/sandwizz-promises-to-reinvent-the-breadboard/",
		"https://hackaday.com/2024/05/26/2024-business-card-challenge-adding-some-refinement-to-breadboard-power-supplies/",
		"https://www.cnx-software.com/2024/05/21/breadboardos-firmware-for-the-raspberry-pi-rp2040-features-a-linux-like-terminal/",
		"https://www.theverge.com/24152153/animal-well-review-switch-ps5-steam-videogamedunkey",
		"https://www.techspot.com/news/103113-snapdragon-x-windows-pcs-run-over-1000-games.html",
	})
	fmt.Println("### RETRIEVAL ###")
	datautils.ForEach(beansack.Retrieve(options), func(item *beansack.Bean) {
		fmt.Printf("[%s] Text length = %d: %s\n", item.Source, len(item.Text), item.Title)
	})

	// trending nuggets
	nuggets := beansack.TrendingNuggets(beansack.NewSearchOptions().WithTimeWindow(2))
	log.Println("### TRENDING NUGGETS ###")
	datautils.ForEach(nuggets, func(item *beansack.BeanNugget) {
		fmt.Printf("%d | %s -> %s | urls = %d\n", item.TrendScore, item.KeyPhrase, item.Event, len(item.BeanUrls))
	})
}

func Search() {
	// initialize the services
	if err := beansack.InitializeBeanSack(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("EMBEDDER_URL"), getEmbedderCtx(), os.Getenv("LLMSERVICE_API_KEY")); err != nil {
		log.Fatalln("initialization not working", err)
	}

	// search and query
	var beans []beansack.Bean
	search_texts := []string{"Cybersecurity", "Security Breach", "Threat Intelligence"}

	// test vector search
	search_opt := beansack.NewSearchOptions().WithTopN(10).WithTimeWindow(3)
	search_opt.SearchTexts = search_texts
	beans = beansack.FuzzySearch(search_opt)
	log.Println("### Category Search Result ###")
	datautils.ForEach(beans, func(item *beansack.Bean) { log.Printf("%f | %s\n", item.SearchScore, item.Title) })

	// test context search
	search_opt = beansack.NewSearchOptions().WithTopN(10).WithTimeWindow(2)
	search_opt.Context = search_texts[0]
	beans = beansack.FuzzySearch(search_opt)
	log.Println("### Context Search Result ###")
	datautils.ForEach(beans, func(item *beansack.Bean) { log.Printf("%f | %s\n", item.SearchScore, item.Title) })

	// test text search
	search_opt = beansack.NewSearchOptions().WithTopN(10).WithTimeWindow(2)
	beans = beansack.TextSearch(search_texts, search_opt)
	log.Println("### Text Search Result ###")
	datautils.ForEach(beans, func(item *beansack.Bean) { log.Printf("%f | %s\n", item.SearchScore, item.Title) })

	// nugget search
	beans = beansack.NuggetSearch(search_texts, beansack.NewSearchOptions().WithTimeWindow(2))
	log.Println("### Nuggets Search Result ###")
	datautils.ForEach(beans, func(item *beansack.Bean) {
		log.Printf("[%s] %s | %s\n", item.Source, time.Unix(item.Updated, 0).Format(time.DateTime), item.Title)
	})
}

func Nlp() {
	// inputs := []string{
	// 	"Johnson & Johnson (NYSE: $JNJ) is an American multinational pharmaceutical and medical technologies corporation headquartered in New Brunswick, New Jersey, founded by three brothers in 1886. The company owns many household names in the healthcare consumer products sectors.  It reported its first-quarter financial results for fiscal 2024 on Tuesday, April 16, 2024, showcasing a mixed performance amidst a dynamic healthcare landscape.  Johnson & Johnson reported a strong performance in the first quarter of 2024, showcasing its ability to navigate the dynamic healthcare landscape. The pharma giant reported adjusted earnings per share (EPS) of $2.71, a 12.4% increase year-over-year, surpassing the consensus estimate of $2.64.  Sales reached $21.38 billion, a 2.3% year-over-year increase, narrowly missing the consensus of $21.39 billion. Operational growth stood at 3.9%, with adjusted operational growth reaching 4.0%. Notably, it reported a net profit of $5.35 billion, a significant improvement from the net loss of $491 million in the same period last year, which included $6.9 billion in litigation charges. The Innovative Medicine segment was a standout performer, with worldwide operational sales, excluding the COVID-19 vaccine, growing 8.3% to $13.6 billion. Key products like the psoriasis drug Stelara and the cancer treatment Darzalex contributed to this strong performance. Additionally, the company’s medical devices business saw a 4.5% year-over-year increase, driven by electrophysiology, cardiovascular, and general surgery product growth. The company’s worldwide sales increased by 3.9%, driven by a notable 7.8% growth in the U.S. market. This top-line growth was accompanied by improvements in profitability, as the company’s adjusted net earnings and adjusted diluted earnings per share rose by 3.8% and 12.4%, respectively, compared to the same period in 2023. Recognizing the company’s financial strength, Johnson & Johnson also announced a quarterly dividend increase of 4.2%, from $1.19 to $1.24 per share.  JNJ 2024 Guidance Despite the solid performance in the first quarter, Johnson & Johnson faced some challenges, particularly in its international markets and within specific product lines. Sales outside the U.S. declined by 0.3%, indicating headwinds in certain global regions. The company’s Contact Lens business also saw a 2.3% decline, driven by U.S. stocking dynamics. To address these challenges, the company has remained focused on its strategic priorities, including strengthening its pipeline, driving commercial excellence, and optimizing its operations. Recent regulatory and clinical milestones, such as FDA approvals and positive study results, have bolstered the company’s confidence in its ability to maintain its growth trajectory. Looking ahead, the company has narrowed its full-year 2024 guidance. The company now forecasts total sales between $88 billion and $88.4 billion, compared to the previous range of $87.8 billion to $88.6 billion. The company’s adjusted EPS guidance has also been adjusted to $10.57 to $10.72, down from the previous range of $10.55 to $10.75. Despite the slightly revised guidance, management remains confident in its ability to deliver above-market growth and continue creating value for its shareholders. The company is focused on executing its strategic priorities, leveraging its diverse portfolio, and investing in innovative solutions. Johnson & Johnson Segment Performance The Innovative Medicine segment reported a 2.5% worldwide sales increase, driven by the strong performance of key brands and the successful uptake of recently launched products. Notable contributors included the company’s immunology and oncology portfolios, with standout products such as Stelara and Darzalex showcasing solid growth. The MedTech segment delivered an impressive 6.3% sales increase globally, with strong growth in both the U.S. and international markets. This performance was supported by the company’s electrophysiology products and the Abiomed cardiovascular business and robust growth in wound closure products within the General Surgery division. JNJ Stock Performance Johnson & Johnson’s stock (NYSE: $JNJ) traded down by over 2.13% on the day following the release of its first-quarter 2024 earnings report. The stock closed at $144.45, a decline of $3.14 compared to the previous trading day. On Wednesday, April 17, the stock recouped some of the previous day’s losses, rising by 0.24% to $144.79 as of 9:33 AM EDT.  Tuesday’s decline was steeper than the 0.2% decline of the broader S&P 500 index, indicating that investors were not entirely satisfied with the company’s mixed quarterly results. The lowered top end of Johnson & Johnson’s full-year guidance range likely contributed to the stock’s sell-off, as the market reacted cautiously to the company’s revised outlook.  Johnson & Johnson (NYSE: $JNJ) Should You Consider Buying Johnson & Johnson Stock in 2024? While Johnson & Johnson’s first-quarter results showed mixed performance, it remains a diversified healthcare leader with a strong portfolio and pipeline. Investors should thoroughly evaluate the company’s long-term growth potential and the risks and opportunities before deciding whether to invest in JNJ stock in 2024. A balanced assessment of all available information is advisable when making investment decisions. Click Here for Updates on Johnson & Johnson – It’s 100% FREE to Sign Up for Text Message Notifications!  Disclaimer: This website provides information about cryptocurrency and stock market investments. This website does not provide investment advice and should not be used as a replacement for investment advice from a qualified professional. This website is for educational and informational purposes only. The owner of this website is not a registered investment advisor and does not offer investment advice. You, the reader / viewer, bear responsibility for your own investment decisions and should seek the advice of a qualified securities professional before making any investment.",
	// 	"Google is adding a bunch of new features to its Gemini AI, and one of the most powerful is a personalization option called “Gems” that allows users to create custom versions of the Gemini assistant with varying personalities.Gems lets you create iterations of chatbots that can help you with certain tasks and retain specific characteristics, kind of like making your own bot in Character.AI, the service that lets you talk to virtualized versions of popular characters and celebrities or even a fake psychiatrist. Google says you can make Gemini your gym buddy, sous-chef, coding partner, creative writing guide, or anything you can dream up. Gems feels similar to OpenAI’s GPT Store that lets you make customized ChatGPT chatbots.You can set up a gem by telling Gemini what to do and how to respond. For instance, you can tell it to be your running coach, provide you with a daily run schedule, and to sound upbeat and motivating. Then, in one click, Gemini will make a gem for you as you’ve described. The Gems feature is available “soon” to Gemini Advanced subscribers. Related:",
	// 	"Skip to content\n\n    \n        \n\n        \n            \n\n    \n\n    \n        You might think that jailbreaking a PS4 to run unsigned code is a complicated process that takes fancy tools and lots of work. While developing said jailbreaks was naturally no mean feat, thankfully they’re far easier for the end user to perform. These days, all you need is an LG TV.\nOf course, you can’t just use any LG TV. You’ll need a modern LG webOS smart TV, and you’ll need to jailbreak it before it can in turn be used to modify your PS4. Once that’s done, though, you can install the PPLGPwn tool for jailbreaking PS4s. It’s based on the PPPwn exploit released by [TheFlow], which was then optimized by [xfangxfang] and implemented for LG Smart TVs by [zauceee]. Once installed, you just need to hook up your PS4 to the TV via the Ethernet port. Then, with the exploit running on the TV, telling the PS4 to set up the LAN via PPPoE will be enough to complete the jailbreak.\nThere are other ways to jailbreak a PS4 that don’t involve the use of a specific television. Nonetheless, it’s neat to see the hack done in such an amusing way.\n\n\nThanks to [eyeoncomputers] for the tip!",
	// 	"The classic competitive game show Jeopardy! is getting a spinoff series on Prime Video. Sony Pictures Television is producing Pop Culture Jeopardy!, which turns the classic academic quiz show with three challengers into a team-based trivia game that touches on music, movies, culture, sports, celebrities, entertainment, and more.The new series marks the first time the company’s game show division is expanding the Jeopardy! brand to a streaming platform with a new show. Years ago, the company let watchers binge the main show on Hulu.While the series sounds like a casual side mission for the franchise, Sony Pictures Television’s president of game shows, Suzanne Prete, states it will “be a nail biter” for fans and that teams will compete “at the highest level.” Pop Culture Jeopardy! is backed by producer Michael Davies, who previously worked on Who Wants to Be a Millionaire. The show was a huge success when it premiered in the US — giving this new Jeopardy! series some promise.The mainline Jeopardy! series has had several guest hosts since Alex Trebek’s death in 2020 but is now finally settled on Ken Jennings. A host for Pop Culture Jeopardy! has yet to be announced, but here’s hoping the process isn’t quite as long or dramatic.",
	// 	"This is what it took to make a small city safe.In 1992, East Palo Alto was dubbed the “murder capital” of the U.S., with 42 murders in its 2.5 square miles — a per capita rate   higher than that of any other city of any size. In 2023, according to East Palo Alto Police Department statistics released last week, the turnaround seemed complete: zero homicides.Law enforcement leaders, residents and city officials point to a complicated mix of circumstances that turned a crime-ridden community into what the mayor now calls “one of the safest places to live in the peninsula.”The San Francisco Peninsula that Mayor Antonio López referred to is home to Stanford University, the opulent town of Atherton and well-heeled Palo Alto. Residents and city leaders scoff at the overly simple idea that gentrification solved  the city’s problems, although the median household income has drastically increased, and the typical home price is a little more than $900,000.They argue that poverty and crime don’t necessarily go hand in hand. They point to increased development since they earned  the grim title of murder capital, including an Ikea and a Four Seasons hotel. Also: more job opportunities, programs for youth and community policing. And time.“In spite of the wrongs of our past, we can move forward and be a model for everyone,” López said.::                     Community Service Officer Magd Algaheim giving a young community member a sticker during a National Night Out party at Joel Davis Park in 2023.   (Jeff Liu/East Palo Alto Police Dept.)       East Palo Alto, wedged between San Francisco Bay and Palo Alto, was home to  a little more than 28,000 people in the last census count. The median household income is $103,000. As tech companies in surrounding areas flourished, the city experienced significant gentrification.  Rising housing prices left some low-income East Palo Altans feeling pushed out. Now majority Latino, East Palo Alto was once a landing spot for Black families facing  redlining by lenders and real estate agents. It’s where Paul Bains’ family landed in 1962, at a time when the city  was majority Black.“East Palo Alto was considered like a small Black mecca,” said Bains, who lives and works there to this day.Many people of color owned  homes and looked out for one another, he said. Bains,  the pastor for St. Samuel Church, refers to that time  as “BC” — before crack.During the nation’s crack cocaine epidemic, about 17% of East Palo Alto residents lived in poverty, higher than the national level.“When crack came along, it just demoralized families,” Bains said.East Palo Alto was also  enduring growing pains. After years of neglect by the San Mateo County government, the community voted to incorporate as a city in 1983. Paul Norris started with the newly formed East Palo Alto Police Department  soon after, in 1987, when there were “drug sales on almost every street.” Although there were only four homicides that year, the number increased dramatically in the years that followed, hitting double digits.                     East Palo Alto Sergeant Matafanua Lualemaga interacts with community members and a miniature horse at the Azalia Drive block party during National Night Out in 2022.   (Jeff Liu/East Palo Alto Police Dept.)       The majority of homicides, Norris said, were tied to drug dealers and gang members fighting over territory.Norris,  now an acting sergeant with the department,  recalled working a night shift when there were six shootings around the city. He waited with a body for 45 minutes before paramedics could arrive.“It just seemed like a war zone,” he said.  Clyde Virges was on the front lines. Raised in East Palo Alto, he had gone away  to college and returned in the late 1980s with a drug problem. Drugs were so plentiful, he said, he could find them littering the street, the fallen crumbs of sales to users who had driven to the troubled city to buy.“You’d take a match cover and scrape up what you could get off the ground,” he said.Gunfire was the city’s soundtrack. Residents used to say that in Palo Alto, there were drive-ins;  in East Palo Alto, there were drive-bys.One New Year’s Eve, Virges  pulled a flower pot over his head  to protect himself as he went to buy crack cocaine, he recalled.Virges found himself caught up in the crackdown in the ’90s, getting arrested after selling controlled substances to an undercover informant who caught him on video. He called it a blessing. He spent nearly a year in a recovery home before going to school to become a licensed drug and alcohol counselor. Now 70, he works as a case manager with the nonprofit organization  WeHope in East Palo Alto, helping the homeless.::In 1992, East Palo Alto gained nationwide notoriety.At the time, the city was home to 24,000 people, and it recorded the highest murder rate in the U.S. Homicides jumped to 42 from 20 the year before.                     The body of a shooting victim lies just behind gang graffiti in an alley off Fordham Street in East Palo Alto in 2000.    (Carlos Avila Gonzalez / San Francisco Chronicle via Getty Images)        It was the equivalent of 175 slayings for every 100,000 residents. “It brought us a hell of a lot of attention here in the state of California,” said Burnham Matthews, the  police chief at the time. “It was a city that really needed some help.”Help came in the form of outside law enforcement. Palo Alto donated four officers; Menlo Park provided two. Later, the San Mateo County Sheriff’s Department brought in 18 deputies, and the state assigned 12 Highway Patrol officers.  The outside officers more than doubled the strength of the department.By 1993, according to police statistics, the number of homicides dropped to four.Sharifa Wilson, who was mayor in 1992, said she welcomed the police support at the time but  stressed that it was “not the answer.”“Part of the issue ... was a lack of economic opportunity,” the 73-year-old said. “We didn’t have access to capital to help establish ourselves.”Thanks to the attention on East Palo Alto, she said, the city was able to get help from the state and move forward with the development of a shopping center — with a policy in place requiring all businesses to hire from the community. “We don’t raise our kids to be drug dealers,” she said. “By creating opportunities for them to work, that had an impact.”Wilson, who still lives in the city, also  credited the community for the reduction in crime. At one point, residents formed a group called “Just Us” and would go to  street corners and take pictures of the  license plates of those driving in  to buy drugs. From there, police sent letters to the registered owners  notifying them that their cars had been seen in a high-drug, high-crime area. (Wilson said one of the letters went to a judge whose son was using the car to buy drugs.)Local nonprofit and faith-based groups  focused on engaging youth in after-school programs and activities that would steer them away from crime.“It really is a testament to the commitment of the community to fix itself,” Wilson said. “East Palo Alto has always been a resilient community. People there are really concerned and care about the community where they live.“The fact that we were labeled the homicide capital gave us an attention that we needed, and then we took that attention and turned it into something positive,” she added. “If you give us lemons, we’re gonna make lemonade.”::                        (East Palo Alto Police Department)       For 17 years, the number of homicides in East Palo Alto remained in the single digits. In 2017 and again in 2019, there was one murder in the city.Shortly after midnight this  New Year’s Day, Police Chief Jeff Liu texted city officials with the news they had long hoped to hear. Finally, the number of homicides in East Palo Alto had fallen to zero.“We’ve always had at least one, and to reach zero is just such a monumental achievement for our whole community,” Liu said. “It’s like the goal that always slipped through our fingers.”Along with the work of his police force, Liu  credited  residents and the efforts they put into reducing crime,  through helping youth and alerting officers to what was happening in the city. That wouldn’t be possible, he said, if the department hadn’t built up trust.When  López heard the news, he said, “it almost made me cry.” The mayor  cited the investment the city made in  funding for the police force so its  salaries would be on par with those of other law enforcement agencies.“What I love about East Palo Alto is not only is it a model for our peninsula but also the country about how community policing can be effective,” he said.With  homicides  down to zero, leaders are  optimistic about  the future.  Bains, the pastor and co-founder and president of WeHope,  said he’s “proud of our city.”“There’s so many heroes that we stand on the shoulders of that stayed in this community and saw the best of times before, to the worst of times, now back to the better of times,” he said. “Now that we have zero murders, we want to keep it at zero murders.”    More to Read",
	// 	"There are a lot of hacking gadgets that can be used to pen test stuff. Like a bad usb, Flipper Zero, deauther watch, pwnagotchi, etc etc. But couldn't I just use my Laptop for those kinds of things? Hardware wise its probably better than those gadgets.\nIm new to pen testing and was just wondering if one just couldn't use their laptop to do the same stuff that those gadgets can.",
	// 	"about this project\nI am passionate about self development and programming. Thus I created an affirmations app with Go - which is becoming quickly my favourite language.\nThe code is by no means perfect, although I tried to keep it as clean as possible.\nI've used Echo and Templ on backend, materialize.css and jQuery on frontend.\nFeel free to modify it and improve it as you see fit.\nYou can see it in action and use it here:\nhttps://easyaffirm.com/\ntooling needed to be installed before running\napt-get install xdotool\napt install mariadb-server\napt-get install npm\nnpm install -g sass\nnpm install -g browserify\ncd internal/frontend\nnpm install\ngo mod tidy\ngo install github.com/cosmtrek/air@latest\ngo install github.com/a-h/templ/cmd/templ@latest\nalso install the migrate tool.\nhttps://github.com/golang-migrate/migrate/releases\nYou can download the .deb package from assets and double click to install\nmysql setup:\nthen:\nCREATE DATABASE affirmtempl;\nCREATE USER 'affirmtempl'@'localhost';\nGRANT ALL PRIVILEGES ON affirmtempl.* TO 'affirmtempl'@'localhost' WITH GRANT OPTION;\n-- Important: Make sure to swap 'pass' with a password of your own choosing.\nALTER USER 'affirmtempl'@'localhost' IDENTIFIED BY 'pass';\nthen migrate\nmigrate -path=./migrations -database=\"mysql://affirmtempl:pass@tcp(localhost)/affirmtempl\" up\nrun dev\nrun air in the root directory of this project to run a development server and live reload.\nOpen 'localhost:4000'  in your browser to use this app.\nthe reload script will try to find a Google Chrome window for hot reloading\nin the script located at ./scripts/devreload.sh\nWID=$(xdotool search --name \"Google Chrome\")\nIf you have Firefox or another browser, change the name to the respective browser instead of \"Google Chrome\"\nstyle modifications\nstatic/sass/app/_appstyle.scss\nstatic/sass/materialize.scss\nstatic/sass/_color-variables.scss\nstatic/sass/components/_variables.scss\nstatic/js/app.js\nlicensing and disclaimer\nI'm releasing this under a Creative Common CC0 license. Basically you can do whatever you want with it, no need for attribution either.\nhttps://creativecommons.org/public-domain/cc0/\nDisclaimer:\nTHE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.",
	// }
	inputs := datautils.Transform(getBeans("./examples/testdata/dataset2.json"), func(item *beansack.Bean) string { return nlp.TruncateTextOnTokenCount(item.Text) })

	// for embeddings
	embed := nlp.NewEmbeddingsDriver(os.Getenv("EMBEDDER_URL"), getEmbedderCtx())
	start_time := time.Now()
	res := datautils.ForEach(embed.CreateBatchTextEmbeddings(inputs, nlp.SEARCH_DOCUMENT), func(emb *[]float32) {
		if len(*emb) == 0 {
			log.Println("returned dud")
		} else {
			log.Println((*emb)[0])
		}
	})
	dur := time.Since(start_time) / time.Second
	log.Printf("%d embeddings generated in %ds. Avg %f\n", len(res), dur, float32(dur)/float32(len(res)))

	// for keyconcepts and digests
	pb := nlp.NewParrotboxClient(os.Getenv("LLMSERVICE_API_KEY"))

	digests := pb.ExtractDigests(inputs)
	fmt.Println(datautils.ToJsonString(digests))

	nuggets := pb.ExtractKeyConcepts(inputs)
	fmt.Println(datautils.ToJsonString(nuggets))
}

func NewBeans() {
	beans := getBeans("./examples/testdata/dataset1.json")
	// initialize the services
	if err := beansack.InitializeBeanSack(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("EMBEDDER_URL"), getEmbedderCtx(), os.Getenv("LLMSERVICE_API_KEY")); err != nil {
		log.Fatalln("Beansack initialization not working.", err)
	}
	log.Println(len(beans), "New Beans")
	beansack.AddBeans(beans)
	// add it again for testing medianoise updates
	beansack.AddBeans(beans)
}

func getBeans(dataset string) []beansack.Bean {
	file, err := os.Open(dataset)
	if err != nil {
		log.Fatalln("Error opening dataset.", err)
	}
	var beans []beansack.Bean
	if err = json.NewDecoder(file).Decode(&beans); err != nil {
		log.Fatalln("Error decoding file.", err)
	}
	return beans
}

func getEmbedderCtx() int {
	ctx, _ := strconv.Atoi(os.Getenv("EMBEDDER_CTX"))
	return ctx
}

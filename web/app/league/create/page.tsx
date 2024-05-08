export default function LeagueCreatePage(){
    const options: string[] = ["1","2","3","4","5"]
    return (
        <div>
            <form>
                <div>
                    <input placeholder="League Name"/>
                </div>
                <div>
                    <label htmlFor="underline_select" className="sr-only">Underline select</label>
                    <select id="underline_select"
                            className="block py-2.5 px-0 w-full text-sm text-gray-500 bg-transparent border-0 border-b-2 border-gray-200 appearance-none dark:text-gray-400 dark:border-gray-700 focus:outline-none focus:ring-0 focus:border-gray-200 peer">
                        <option selected>Choose the level</option>
                        <option value="1">1</option>
                        <option value="2">2</option>
                        <option value="3">3</option>
                        <option value="4">4</option>
                        <option value="5">5</option>
                    </select>
                </div>
            </form>
        </div>
    )
}
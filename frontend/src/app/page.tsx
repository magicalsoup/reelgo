import Navbar from "@/components/nav/Navbar";

// landing
export default async function Home() {
  return (
    <div className="h-screen">
      <Navbar/>
      <div className="flex flex-col h-full justify-center items-center">
        <div className="flex flex-row">
          <h1 className="text-9xl font-sans w-[820px]">Turning your reels to a travel itinerary</h1>
          <img width="400" height="400" className="bg-gray-200"/>
        </div>
      </div>
    </div>
  );
}
